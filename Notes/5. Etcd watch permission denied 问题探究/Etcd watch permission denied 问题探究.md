## 前言
etcd使用watch出现permission denied的报错，找到类似issue：[github issue](https://github.com/etcd-io/etcd/issues/12385)。初步判断为token过期导致的问题，这里从watch的调用链和token校验的逻辑两方面来看。本文代码为etcd master分支4f0e92d。

### Watch调用链

1. Watch函数

```go
// etcd/client/v3/watch.go:290

// Watch posts a watch request to run() and waits for a new watcher channelfunc 
(w *watcher) Watch(ctx context.Context, key string, opts ...OpOption) WatchChan {
	...
	
	// 请求参数
	wr := &watchRequest{  
		...
		retc:           make(chan chan WatchResponse, 1),  // 回调chan
	}

	...

	wgs := w.streams[ctxKey]  // grpc stram对象，实际执行请求

	...
	
	reqc := wgs.reqc // 请求chan

	...
}
```

2. watchGrpcStream run函数内实际执行Send位置

```go
// etcd/client/v3/watch.go:496

// run is the root of the goroutines for managing a watcher client
func (w *watchGrpcStream) run() {
    for {  
        select {  
        // Watch() requested  
        case req := <-w.reqc:  
            switch wreq := req.(type) {  
            case *watchRequest:
            
            ...

            // head of resume queue, can register a new watcher  
            if err := wc.Send(ws.initReq.toPB()); err != nil {  
                w.lg.Debug("error when sending request", zap.Error(err))  
            }  

            ...
}
```

3. Recv位置

```go
// etcd/server/etcdserver/api/v3rpc/watch.go:239

func (sws *serverWatchStream) recvLoop() error {  
    for {  
        req, err := sws.gRPCStream.Recv()

        ...
    }

    // 这里是返回permiss denied的地方
    if !sws.isWatchPermitted(creq) {  
        wr := &pb.WatchResponse{  
        Header:       sws.newResponseHeader(sws.watchStream.Rev()),  
        WatchId:      creq.WatchId,  
        Canceled:     true,  
        Created:      true,  
        CancelReason: rpctypes.ErrGRPCPermissionDenied.Error(),  // 返回的错误
    }

    ...
}
```

注意这里会导致另一个bug，后面细说，跳过一些函数，直接找到token校验的地方

4. token检查位置

etcd/server/auth/store.go:1034

```go
func (as *authStore) AuthInfoFromCtx(ctx context.Context) (*AuthInfo, error) {
    ...

    // 这里校验token失败
    authInfo, uok := as.authInfoFromToken(ctx, token)  
    if !uok {  
        as.lg.Warn("invalid auth token", zap.String("token", token))  
        return nil, ErrInvalidAuthToken  
    }
}
```

这里的日志`invalid auth token ...`正好有打印，所以锁定了位置，下面看token校验的逻辑。

### Token校验逻辑
1. Token在服务器如何校验

```go
// etcd/server/auth/store.go:186

type TokenProvider interface {  
    info(ctx context.Context, token string, revision uint64) (*AuthInfo, bool)  
    assign(ctx context.Context, username string, revision uint64) (string, error)  
    enable()  
    disable()  
  
    invalidateUser(string)  
    genTokenPrefix() (string, error)  
}
```

服务器这里定义了一个接口实现了token的逻辑，内部分别有`tokenJWT`，`tokenNop`，`tokenSimple`实现了接口，默认是使用`tokenSimple`，这里我们也是使用默认设置。

2. tokenSimple的主要逻辑

token校验失败主要是因为token是有过期时间的

```go
// etcd/server/auth/simple_token.go:186

func (tm *simpleTokenTTLKeeper) run() {
    ...

    // 一秒的ticker
    case <-tokenTicker.C:  
        nowtime := time.Now()  
        tm.mu.Lock()  
        for t, tokenendtime := range tm.tokens {  
            // 超时的token都删掉
            if nowtime.After(tokenendtime) {  
                tm.deleteTokenFunc(t)  
                delete(tm.tokens, t)  
	        }  
        }  
        tm.mu.Unlock()
}
```

所以如果我们启动etcd使用以下参数
`--log-level debug --auth-token-ttl 5`
token过期时间为5s，会看到这样的日志：
`{"level":"info","ts":"2022-08-01T16:30:08.906+0800","caller":"auth/simple_token.go:165","msg":"deleted a simple token","user-name":"root","token":"lUKYuvCtkwWqlYQo.197"}`

那么如果token过期就校验不过，为何只有watch会出现这样的问题？这里再回到客户端去看token刷新的逻辑。

### Token刷新逻辑

1. 从client创建的地方开始

```go
// etcd/client/v3/client.go:385

if cfg.Username != "" && cfg.Password != "" {  
    client.Username = cfg.Username  
    client.Password = cfg.Password  
    client.authTokenBundle = credentials.NewBundle(credentials.Config{})  
}
```

如果设置了用户名和密码，防止频繁校验，就自动启动了token机制，`Bundle`是GRPC的鉴权接口，这里关注`PerRPCCredentials`,

```go
// etcd/client/v3/credentials/credentials.go:110

func (rc *perRPCCredential) GetRequestMetadata(ctx context.Context, s ...string) (map[string]string, error) {  
    rc.authTokenMu.RLock()  
    authToken := rc.authToken  
    rc.authTokenMu.RUnlock()  
    if authToken == "" {  
        return nil, nil  
    }  
    return map[string]string{rpctypes.TokenFieldNameGRPC: authToken}, nil  
}
```

这里负责把token注入context，这样就可以把刷新token传到服务端。

2. 如何刷新token

通过Authenticate获取新的token，即通过username和password获取token

```go
// Authenticate login and get token
Authenticate(ctx context.Context, name string, password string) (*AuthenticateResponse, error)
```

3. 调用Authenticate的地方

```go
// etcd/client/v3/client.go:259

func (c *Client) getToken(ctx context.Context) error {
	...

    // Authenticate请求获取新的token
    resp, err := c.Auth.Authenticate(ctx, c.Username, c.Password)  
    if err != nil {  
        if err == rpctypes.ErrAuthNotEnabled {  
            return nil  
        }  
        return err  
    }

    ...
}
```

4. 拦截器调用getToken

```go
// etcd/client/v3/retry_interceptor.go:39

func (c *Client) unaryClientInterceptor(optFuncs ...retryOption) grpc.UnaryClientInterceptor 

// etcd/client/v3/retry_interceptor.go:113

func (c *Client) streamClientInterceptor(optFuncs ...retryOption) grpc.StreamClientInterceptor 
```

5. 总结
一元GRPC请求每次通过`UnaryClientInterceptor`和`PerRPCCredentials`获取并且刷新token保证token的正确性，如果token过期会进行重试。

```go
// etcd/client/v3/retry_interceptor.go:156

// shouldRefreshToken checks whether there's a need to refresh the token based on the error and callOptions,// and returns a boolean value.  
func (c *Client) shouldRefreshToken(err error, callOpts *options) bool {  
    if rpctypes.Error(err) == rpctypes.ErrUserEmpty {  
        // refresh the token when username, password is present but the server returns ErrUserEmpty  
        // which is possible when the client token is cleared somehow      return c.authTokenBundle != nil // equal to c.Username != "" && c.Password != ""  
    }  
  
    return callOpts.retryAuth &&  
        (rpctypes.Error(err) == rpctypes.ErrInvalidAuthToken || rpctypes.Error(err) == rpctypes.ErrAuthOldRevision)  
}
```

这里看到返回`ErrInvalidAuthToken`和`ErrAuthOldRevision`都会刷新Token。

但是对于`GRPCStream`，这里不论是是`StreamClientInterceptor`还是`PerRPCCredentials`都只有在创建的时候调用一次，因为后续watch公用一个`GRPCStream`，所以token不会进行刷新，后续的watch便会失败，目前issue还在打开状态，仍未解决。[my comment](https://github.com/etcd-io/etcd/issues/12385#issuecomment-1203583417)

### 另一个坑

```go
// etcd/server/etcdserver/api/v3rpc/watch.go:239

func (sws *serverWatchStream) recvLoop() error {  
    for {  
        req, err := sws.gRPCStream.Recv()

        ...
    }

    // 这里是返回permiss denied的地方
    if !sws.isWatchPermitted(creq) {  
        wr := &pb.WatchResponse{  
        Header:       sws.newResponseHeader(sws.watchStream.Rev()),  
        WatchId:      creq.WatchId, // 默认=0
        Canceled:     true, // 注意这里
        Created:      true,  
        CancelReason: rpctypes.ErrGRPCPermissionDenied.Error(),  // 返回的错误
    }

    select {
        case sws.ctrlStream <- wr: // 发送到ctrlStream
            continue
        case <-sws.closec:
            return nil
    }

    ...

    case c, ok := <-sws.ctrlStream: // 从ctrlStream获取
    if !ok {
        return
    }

    ...

    // track id creation
    wid := mvcc.WatchID(c.WatchId)
    if c.Canceled { // Canceled==true
        delete(ids, wid) // 从ids删除，ids是活跃的WatchID，客户端WatchID=0的watch再也收不到消息了
        continue
    }

	...
}
```

这里会导致另一个坑，注意`WatchID=0`和`Canceled=true`两个条件会导致`WatchID=0`从活跃的`WatchID`中删除，导致客户端`WatchID=0`的watch再也接受不到信息。这里考虑把初始`WatchID`设置为1解决bug。[fix PR](https://github.com/etcd-io/etcd/pull/14296)

### 后记

顺便提一下watch的坑，etcd的revision机制，注意看watch的这一段注释：
```go
// If revisions waiting to be sent over the  
// watch are compacted, then the watch will be canceled by the server, the  
// client will post a compacted error watch response, and the channel will close.
```


