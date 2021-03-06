## Rocky Linux 简介
[rocky linux官方简介](https://rockylinux.org/about)

Red Hat宣布停止维护RHEL下游版本CentOS，转而支持上游名为"CentOS Stream"的版本。为了寻回初心Centos原先的创始人Gregory Kurtzer宣布推出Rocky Linux。

## 安装Rocky Linux 8
1. 首先下载镜像，一般选择Minimal或者Boot镜像，看下方描述是一样的。[rocky linux官方镜像下载](https://rockylinux.org/download)
2. 选择下载的镜像，进行安装。
3. 安装过程中设置完红色字体的模块后完成安装。

## 设置静态IP
目的是虚拟机重启之后ip不会变化，方便自己使用。

#### 方法一(推荐)：
1. 打开图形化设置界面
```
nmtui
```

2. 如图设置 ![](https://github.com/kafuu-chino/learning-notes/blob/main/Notes/1.%20VMware%E5%AE%89%E8%A3%85Rocky%20Linux%208/imgs/1.1.png?raw=true)

#### 方法二：
1. 打开配置文件：
```
vim /etc/sysconfig/network-scripts/ifcfg-ens160
```

2. 修改配置文件：
```
BOOTPROTO=static
IPADDR=192.168.178.130
PREFIX=24
GATEWAY=192.168.178.2
DNS1=192.168.178.2
```

以上方法设置完毕后，重启生效。