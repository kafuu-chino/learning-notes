## 简介
Apollo是携程开发的分布式配置中心，方便不同环境的配置管理，配置修改实时生效等优点。更多参考[官方介绍](https://www.apolloconfig.com/#/zh/README)。

对比阿里开发的Nacos，Apollo对于权限功能更加全面，而且稳定性更好。[更多详情对比](https://www.lmf.life/archives/apollo%E4%B8%8Enacos%E9%85%8D%E7%BD%AE%E4%B8%AD%E5%BF%83%E5%8A%9F%E8%83%BD%E7%9A%84%E5%AF%B9%E6%AF%94)。

## 安装
参考[官方部署教程](https://www.apolloconfig.com/#/zh/deployment/distributed-deployment-guide)。

### 准备工作：
1. 安装java：
```
dnf install -y java-17-openjdk
```

2. 安装mysql，详情请见[mysql官方安装指南](https://dev.mysql.com/doc/refman/8.0/en/installing.html)。

3. 导入两个sql并验证：
[sql文件链接](https://github.com/apolloconfig/apollo/tree/master/scripts/sql)

导入sql：
```
mysql> source /root/apollo/apolloportaldb.sql
mysql> source /root/apollo/apolloconfigdb.sql
```

验证：
```
mysql> select `Id`, `Key`, `Value`, `Comment` from `ApolloPortalDB`.`ServerConfig` limit 1;
+----+--------------------+-------+--------------------------+
| Id | Key                | Value | Comment                  |
+----+--------------------+-------+--------------------------+
|  1 | apollo.portal.envs | dev   | 可支持的环境列表         |
+----+--------------------+-------+--------------------------+

mysql> select `Id`, `Key`, `Value`, `Comment` from `ApolloConfigDB`.`ServerConfig` limit 1;
+----+--------------------+-------------------------------+------------------------------------------------------+
| Id | Key                | Value                         | Comment                                              |
+----+--------------------+-------------------------------+------------------------------------------------------+
|  1 | eureka.service.url | http://localhost:8080/eureka/ | Eureka服务Url，多个service以英文逗号分隔             |
+----+--------------------+-------------------------------+------------------------------------------------------+
```

### 方法一：使用安装包安装
1. releases页面下载压缩包：
[releases页面](https://github.com/apolloconfig/apollo/releases)

```
wget https://github.com/apolloconfig/apollo/releases/download/v2.0.0/apollo-adminservice-2.0.0-github.zip

wget https://github.com/apolloconfig/apollo/releases/download/v2.0.0/apollo-configservice-2.0.0-github.zip

wget https://github.com/apolloconfig/apollo/releases/download/v2.0.0/apollo-portal-2.0.0-github.zip
```

2. 解压：
```
unzip apollo-adminservice-2.0.0-github.zip -d ./adminservice
unzip apollo-configservice-2.0.0-github.zip -d ./configservice
unzip apollo-portal-2.0.0-github.zip -d ./portal
```

3. 启动和停止：
执行解压后的`scripts`文件夹下脚本即可
```
startup.sh
shutdown.sh  
```

### 方法二：使用podman安装
启动项包含部分配置如ip等，需要按需求修改成自己的

1. 启动`portal`
```
podman run --privileged \
    -p 8070:8070 \
    -e SPRING_DATASOURCE_URL="jdbc:mysql://192.168.178.132:3306/ApolloPortalDB?characterEncoding=utf8" \
    -e SPRING_DATASOURCE_USERNAME=root -e SPRING_DATASOURCE_PASSWORD=123456 \
    -e APOLLO_PORTAL_ENVS=dev,pro \
    -e DEV_META=http://192.168.178.128:8080,http://192.168.178.130:8080 \
    -e PRO_META=http://192.168.178.128:8081,http://192.168.178.130:8081 \
    -d -v /root/apollo/portal/logs:/opt/logs --name apollo-portal apolloconfig/apollo-portal
```

参数说明：
* SPRING_DATASOURCE_URL: 对应环境ApolloPortalDB的地址
* SPRING_DATASOURCE_USERNAME: 对应环境ApolloPortalDB的用户名
* SPRING_DATASOURCE_PASSWORD: 对应环境ApolloPortalDB的密码
* APOLLO_PORTAL_ENVS(可选): 对应ApolloPortalDB中的apollo.portal.envs配置项，如果没有在数据库中配置的话，可以通过此环境参数配置
* DEV_META/PRO_META(可选): 配置对应环境的Meta Service地址，以${ENV}_META命名，需要注意的是如果配置了ApolloPortalDB中的apollo.portal.meta.servers配置，则以apollo.portal.meta.servers中的配置为准

2. 启动`adminservice`
```
podman run --privileged \
    -p 8090:8090 \
    -e SPRING_DATASOURCE_URL="jdbc:mysql://192.168.178.132:3306/ApolloConfigDB?characterEncoding=utf8" \
    -e SPRING_DATASOURCE_USERNAME=root -e SPRING_DATASOURCE_PASSWORD=123456 \
    -d -v /root/apollo/adminservice/logs:/opt/logs --name apollo-adminservice apolloconfig/apollo-adminservice
```

3. 启动`configservice`
```
podman run --privileged \
    -p 8080:8080 \
    -e SPRING_DATASOURCE_URL="jdbc:mysql://192.168.178.132:3306/ApolloConfigDB?characterEncoding=utf8" \
    -e SPRING_DATASOURCE_USERNAME=root -e SPRING_DATASOURCE_PASSWORD=123456 \
    -d -v /root/apollo/adminservice/logs:/opt/logs --name apollo-configservice apolloconfig/apollo-configservice
```

## 部署
首先需要了解自己想要部署的架构，详情请看[官网讲解](https://www.apolloconfig.com/#/zh/deployment/deployment-architecture)。

这里部署高可用的dev和pro环境，参考[高可用，双环境](https://www.apolloconfig.com/#/zh/deployment/deployment-architecture?id=_33-%e9%ab%98%e5%8f%af%e7%94%a8%ef%bc%8c%e5%8f%8c%e7%8e%af%e5%a2%83)。
1. 服务器 `192.168.178.132` 安装`mysql`
2. 服务器 `192.168.178.128` 启动`protal`
3. 服务器 `192.168.178.128` 和 `192.168.178.130` 为dev和pro启动两组`adminservice`，`configservice`

如果是安装包安装，`portal`需要设置`conf`下面的`apollo-env.properties`，且三个服务都需要在`conf`下的`application-github.properties`配置数据库信息，这里使用podman进行部署。

首先两套环境需要两个独立的ApolloConfigDB，修改`apolloportaldb.sql`文件，分别改database为ApolloConfigDB_DEV，ApolloConfigDB_PRO。例如dev的sql修改为：
```
CREATE DATABASE IF NOT EXISTS ApolloConfigDB_DEV DEFAULT CHARACTER SET = utf8mb4;

Use ApolloConfigDB_DEV;
```

分别导入mysql后开始部署`adminservice`和`configservice`。

创建好对应的logs文件夹后，两台机器分别执行：
启动dev的`adminservice`：
```
podman run --privileged \
    -p 8090:8090 \
    -e SPRING_DATASOURCE_URL="jdbc:mysql://192.168.178.132:3306/ApolloConfigDB_DEV?characterEncoding=utf8" \
    -e SPRING_DATASOURCE_USERNAME=root -e SPRING_DATASOURCE_PASSWORD=123456 \
    -d -v /root/apollo/dev/adminservice/logs:/opt/logs --name apollo-adminservice_dev apolloconfig/apollo-adminservice
```

启动dev的`configservice`：
```
podman run --privileged \
    -p 8080:8080 \
    -e SPRING_DATASOURCE_URL="jdbc:mysql://192.168.178.132:3306/ApolloConfigDB_DEV?characterEncoding=utf8" \
    -e SPRING_DATASOURCE_USERNAME=root -e SPRING_DATASOURCE_PASSWORD=123456 \
    -d -v /root/apollo/adminservice/logs:/opt/logs --name apollo-configservice_dev apolloconfig/apollo-configservice
```

启动pro的`adminservice`：
```
podman run --privileged \
    -p 8091:8090 \
    -e SPRING_DATASOURCE_URL="jdbc:mysql://192.168.178.132:3306/ApolloConfigDB_PRO?characterEncoding=utf8" \
    -e SPRING_DATASOURCE_USERNAME=root -e SPRING_DATASOURCE_PASSWORD=123456 \
    -d -v /root/apollo/dev/adminservice/logs:/opt/logs --name apollo-adminservice_pro apolloconfig/apollo-adminservice
```

启动pro的`configservice`：
```
podman run --privileged \
    -p 8081:8080 \
    -e SPRING_DATASOURCE_URL="jdbc:mysql://192.168.178.132:3306/ApolloConfigDB_PRO?characterEncoding=utf8" \
    -e SPRING_DATASOURCE_USERNAME=root -e SPRING_DATASOURCE_PASSWORD=123456 \
    -d -v /root/apollo/adminservice/logs:/opt/logs --name apollo-configservice_pro apolloconfig/apollo-configservice
```

查看容器：
```
podman ps

CONTAINER ID  IMAGE                                               COMMAND               CREATED         STATUS             PORTS                   NAMES
234a56573217  docker.io/apolloconfig/apollo-adminservice:latest   /apollo-adminserv...  2 minutes ago   Up 2 minutes ago   0.0.0.0:8090->8090/tcp  apollo-adminservice_dev
7d32d07c81ed  docker.io/apolloconfig/apollo-configservice:latest  /apollo-configser...  2 minutes ago   Up 2 minutes ago   0.0.0.0:8080->8080/tcp  apollo-configservice_dev
a693fe83b10f  docker.io/apolloconfig/apollo-adminservice:latest   /apollo-adminserv...  29 seconds ago  Up 23 seconds ago  0.0.0.0:8091->8090/tcp  apollo-adminservice_pro
82b77300b77b  docker.io/apolloconfig/apollo-configservice:latest  /apollo-configser...  15 seconds ago  Up 8 seconds ago   0.0.0.0:8081->8080/tcp  apollo-configservice_pro
```

注意这里有个问题，如果用podman跑这两个服务，获取本地ip的时候会获取到podman网卡的ip，而不是机器网卡ip，这时候有3种解决方案，具体参考[官方Q&A](https://github.com/apolloconfig/apollo/wiki/%E9%83%A8%E7%BD%B2&%E5%BC%80%E5%8F%91%E9%81%87%E5%88%B0%E7%9A%84%E5%B8%B8%E8%A7%81%E9%97%AE%E9%A2%98#3-admin-server-%E6%88%96%E8%80%85-config-server-%E6%B3%A8%E5%86%8C%E4%BA%86%E5%86%85%E7%BD%91ip%E5%AF%BC%E8%87%B4portal%E6%88%96%E8%80%85client%E8%AE%BF%E9%97%AE%E4%B8%8D%E4%BA%86admin-server%E6%88%96config-server)，这里大概介绍下。
1. 源码修改`application.yml`文件，屏蔽podman网卡，重新打包。
2. 修改`startup.sh`，为JAVA_OPTS添加-Deureka.instance.ip-address=${指定的IP}选项
3. 修改`startup.sh`，为JAVA_OPTS添加-Deureka.instance.homePageUrl=http://${指定的IP}:${指定的Port}选项

这里采用第3种方案：
1. 进入容器：
```
podman exec -it 234a56573217 /bin/bash 
```

2. 打开startup.sh：
```
vi /apollo-adminservice/scripts/startup.sh
```

3. 添加`-Deureka.instance.homePageUrl={ip:port}`，保存退出：
```
export JAVA_OPTS="$JAVA_OPTS -Deureka.instance.homePageUrl=192.168.178.128:8080...
```

4. 重启容器：
```
podman restart 234a56573217
```

重复此操作到其他3个容器。

另一台机器重复以上操作之后，开始部署`portal`。
```
podman run --privileged \
    -p 8070:8070 \
    -e SPRING_DATASOURCE_URL="jdbc:mysql://192.168.178.132:3306/ApolloPortalDB?characterEncoding=utf8" \
    -e SPRING_DATASOURCE_USERNAME=root -e SPRING_DATASOURCE_PASSWORD=123456 \
    -d -v /root/apollo/portal/logs:/opt/logs --name apollo-portal apolloconfig/apollo-portal
```

这里也可以用可选参数配置，这次在mysql配置，更加灵活一些，进入`ApolloPortalDB`的`ServerConfig`表。
1. 修改`apollo.portal.envs`配置环境：
```
pro,dev
```
2. 修改`apollo.portal.meta.servers`为`configservice`服务器的ip端口：
```
{"dev":"http://192.168.178.128:8080,http://192.168.178.130:8080","pro":"http://192.168.178.128:8081,http://192.168.178.130:8081"}
```

接着配置`ApolloConfigDB`的`ServerConfig`表的`eureka.service.url`用于服务发现:
```
http://192.168.178.128:8080/eureka/,http://192.168.178.130:8080/eureka/
```

配置完全部重启一次。

打开管理页面：http://192.168.178.131:8070/

![[3.1.png]]

如图所示，初始账号密码为：apollo,admin

进入系统信息：
![[3.2.png]]

可以看到两个高可用环境已经配置完毕~接下来就是试用阶段。

## 试用
首先明确应用，部门，环境，集群，命名空间等概念。
* 应用顾名思义，就是一个app，应用属于一个部门。
* 部门类似于应用组的概念，管理多个应用。

环境，集群和命名空间主要是区分同一应用配置场景
* 环境比如开发，测试，正式环境的配置区分，依附于应用，可以同时推送。
* 集群如数据中心A和数据中心B的区分，依附于不同环境。
* 命名空间类似于配置文件，公共配置文件可以为其他应用所读取。

这里具体操作请看[官方文档](https://www.apolloconfig.com/#/zh/usage/apollo-user-guide

客户端使用[Agollo - Go Client For Apollo](https://github.com/apolloconfig/agollo)和[viper](https://github.com/spf13/viper)配套使用。

## 总结
实际使用中，代码里面配置的是`configserver`，不可以通过参数直接切换环境，所以推荐走域名负载均衡。生产中，读取哪一组配置也应该有一个类似于配置服务器的去负责，方便管理不同环境和切换。

还有一点体会就是部署略微复杂，而且会有坑需要去修改源码或者去修改启动脚本，无法通过配置文件完成全部操作。

至于使用起来确实会让配置环境更加灵活，多个服务器可以共享一份配置，管理和修改都更方便，而且权限管理可以让出部分工作，比如让QA，策划去修改自己服务器和内部测试服的配置去测试等。导入导出和多个环境推送功能也不会让修改配置变成一个冗余的工作流程。


