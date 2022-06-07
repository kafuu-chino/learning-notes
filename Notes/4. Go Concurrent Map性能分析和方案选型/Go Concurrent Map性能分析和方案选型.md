# Go Concurrent Map性能分析和方案选型
## 前言：
本次讨论并发场景下安全的map，讨论性能问题和使用场景，对map本身实现不过多关注。另附测试代码见文件夹，个人视野有限，有问题欢迎提出。

go中保证map数据并发安全需要加锁，为了提高效率有两种优化方案，
1. 读写分离，使用sync.Map
2. 减少锁的粒度，对map进行分片

所以我们有3种方案，
1. `LockMap`：对map直接加锁
2. `SyncMap`: 读写分离，读不加锁
3. `SliceMap`: 对LockMap进行分片

这里主要针对三种设计之间进行西能对比，其他暂时不做赘述。
按照设计理念，理论上性能应该是：

读操作：`SyncMap` > `SliceMap` > `LockMap`

写操作：`LockMap` > `SliceMap` > `SyncMap`

以下为压测数据。
注意以下两点：
1. 以下压测，都使用不同key，发挥 `SliceMap` 优势，对其他两种map影响不大
2. 涉及到取不同key的操作，使用 `keyPool` ，需要额外消耗30ns/op，压测如下。
```
BenchmarkKeyPool-8   	41412397	        30.30 ns/op	       0 B/op	       0 allocs/op
```

读操作：
```
BenchmarkConcurrenceMap/LockMap_HitGet-8         	20217652	        61.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkConcurrenceMap/SyncMap_HitGet-8         	32106468	        36.93 ns/op	       0 B/op	       0 allocs/op
BenchmarkConcurrenceMap/SliceMap_HitGet-8        	27615411	        45.99 ns/op	       0 B/op	       0 allocs/op
```

写操作：
```
BenchmarkConcurrenceMap/LockMap_MissSet-8        	 7953280	       150.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkConcurrenceMap/SyncMap_MissSet-8        	 5391055	       222.8 ns/op	      32 B/op	       2 allocs/op
BenchmarkConcurrenceMap/SliceMap_MissSet-8       	20545768	        59.82 ns/op	       0 B/op	       0 allocs/op
```

可以看出和猜想一致，注意的是 `SyncMap` 的读操作耗时小于其他两种方案，写入操作大于其他两种方案。对性能要求比较高的场景基本排除 `LockMap` （在简单的场景，对性能要求不高，go还是推荐使用 `LockMap` 的），所以和 `SliceMap` 对比，读多写少的场景选用选用 `SyncMap` ，相反使用 `SliceMap` 。

那么就有个问题，怎么才算读多写少，读和写的比例是多少呢。
这里要了解下`SyncMap`的部分实现原理，内部有缓存数据 `read` 和加锁数据 `dirty`，以下用变量名代替 。
首先了解下不需要加锁的操作：
1. 读操作，命中 `read` 。
2. 读操作，未命中 `read` ，但是没有新增key，这次读取直接判断为miss。
需要加锁的操作：
1. 读操作，未命中 `read` ，有新增key，需要加锁再次读取 `dirty`。
2. 写操作，命中 `read` ，同时修改 `read` 和 `dirty` 。
3. 写操作，未命中 `read` ，修改 `dirty` 。

	总结一下，用 `SyncMap` 写入新key性能较低，读取旧key性能较高，所以判断是否使用 `SyncMap` 的标准可以转化为，写入新key和读取旧key的比例，就是相对 `SliceMap` 写入新key所消耗的时间，要用几次读取旧key补回来，并且其他操作中， `SliceMap` 通过分片，效率也是大于 `SyncMap`  的，这点也需要考虑到。

假设每次写入新key，读取命中的情况，经反复测试，读写比例约2000:1以上， `SyncMap` 才有优势，测数据如下。
```
BenchmarkConcurrenceMap/LockMap_GetMoreThanSet-8         	14390233	        75.92 ns/op	       0 B/op	       0 allocs/op
BenchmarkConcurrenceMap/SyncMap_GetMoreThanSet-8         	30967234	        39.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkConcurrenceMap/SliceMap_GetMoreThanSet-8        	25575064	        46.75 ns/op	       0 B/op	       0 allocs/op
```

这里讨论map上层使用，理论上当map元素越多，虽然hash本身时间复杂度不会增加，但是内部会因为hash冲突等问题，导致性能降低。所以数据量过大的时候 
 `SliceMap` 还有map自身性能的优势，压测如下（依旧使用 `keyPool` ）。
 
 读操作：
 ```
 BenchmarkMap-8-10000       33418690            35.61 ns/op
 BenchmarkMap-8-100000      21556942            59.60 ns/op
 BenchmarkMap-8-1000000     7311390             164.7 ns/op
```

如压测，数量级10w以上读取性能下降就很厉害了。

## 总结：
1. 简单场景，对性能无明显要求使用 `LockMap` 。
2. 读远大于新key写入，数据量不是特别大，使用 `SyncMap` 。
3. 其他场景使用 `SliceMap` 。

## 思考：
理论上 `SliceMap` 也可以用 `sync.Map` 实现，不过个人接触的大部分场景map并不是瓶颈，相反会浪费很多内存，优化要根据实际需要去做。