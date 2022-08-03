安装方法很多，主要是用非root权限并且用包管理工具去安装和启动mysql有一些权限和初始化的坑，详情见[mysql官方安装指南](https://dev.mysql.com/doc/refman/8.0/en/installing.html)。

### 方法一：使用Podman安装(推荐) 
使用podman进行安装，podman是RedHat对于docker的竞品，和docker相比没有守护程序，不需要root权限，且开源，更多参考[podman官方介绍](https://docs.podman.io/en/latest/)。这里可以卸载podman，安装docker，也可以直接使用podman，使用起来大部分情况格式一致，把docker替换成podman使用即可。更多请参考[docker hub安装教程](https://hub.docker.com/_/mysql)。

输入命令：
```
podman run --privileged -d --name mydb -v /root/mysql_data/:/var/lib/mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:8.0.29
```

参数说明：
* --privileged：给予相关权限，否则会报Permission Denied错误，注意这里的格式和docker有些不一样。
* -d：后台运行
* -name：实例名称
* -v：挂载到本地，数据持久化
* -p：映射并暴露端口
* -e：镜像本身启动参数
* MYSQL_ROOT_PASSWORD：root密码
* msyql:8.0.29：mysql镜像版本

如果遇到报错`Can't connect to local MySQL server through socket`，使用tcp connector进行连接，[详情这里](https://stackoverflow.com/questions/4448467/cant-connect-to-local-mysql-server-through-socket-var-lib-mysql-mysql-sock)。
```
mysql -h 127.0.0.1 -u root -p
```

### 方法二：使用包管理工具安装
参考[mysql官方安装教程](https://dev.mysql.com/doc/refman/8.0/en/linux-installation-yum-repo.html)。

#### 拥有root权限：
dnf是作为yum下一代出现的包管理工具，用法基本一致。

1.  使用dnf安装：
```
sudo dnf install mysql-serve
```

2. 启动：
```
sudo mysqld --user=root
```

#### 没有root权限：
pkcon是包管理工具PackageKit的客户端，dnf本质上是root去调用PackageKit，用法依旧没什么区别。

1. 使用pkcon安装：
```
pkcon install mysql-server
```

直接启动会遇到Permission Denied的报错，首先解决权限问题。

查看配置文件/etc/my.cnf.d/mysql-server.cnf：
```
[mysqld]
datadir=/var/lib/mysql
socket=/var/lib/mysql/mysql.sock
log-error=/var/log/mysql/mysqld.log
pid-file=/run/mysqld/mysqld.pid
```

注意`datadir`，`log-error`，`pid-file`都是没有权限的，这里使用用户Sora示范。

2. 创建文件夹`/home/Sora/mysql` 和 `/home/Sora/mysql/data`。

3. 创建配置文件，保存为my.cnf：
```
[mysqld]
datadir=/home/Sora/mysql/data
socket=/var/lib/mysql/mysql.sock
log-error=/home/Sora/mysql/mysqld.log
pid-file=/home/Sora/mysql/mysqld.pid
```

4. 初始化datadir：
```
mysqld --initialize --user=Sora
```

想了解datadir初始化更多，请参考[mysql官方教程](https://dev.mysql.com/doc/mysql-installation-excerpt/5.7/en/data-directory-initialization.html)。

5. 启动mysqld：
```
mysqld --defaults-file=/home/Sora/mysql/my.cnf -D
```

`--defaults-file`指定配置文件，`-D`作为守护进程启动。

6. 从log-error获取root初始密码：
```
grep 'temporary password' /home/Sora/mysql/mysqld.log
A temporary password is generated for root@localhost: w.rsT,2>arD*
```

root用户密码为 `w.rsT,2>arD*`