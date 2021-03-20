# https_expire
zabbix-agent2自定义插件监控https证书过期时间
> 有问题可以联系微信：wanger5354

> 微信公众号：运维开发故事

### 下载zabbix agent2源码并将自定义插件编译
```
yum install golang
git clone https://git.zabbix.com/scm/zbx/zabbix.git --depth 1 zabbix-agent2
cd zabbix-agent2
git submodule add https://github.com/cxf210/ssl_expire.git src/go/plugins/ssl_expire
```
### 导入ssl_expire插件
```
vi src/go/plugins/plugins_linux.go
```
添加最后一行
```go
        _ "zabbix.com/plugins/ceph"
        _ "zabbix.com/plugins/docker"
        _ "zabbix.com/plugins/kernel"
        _ "zabbix.com/plugins/log"
        _ "zabbix.com/plugins/memcached"
        _ "zabbix.com/plugins/modbus"
        _ "zabbix.com/plugins/mqtt"
        _ "zabbix.com/plugins/mysql"
        _ "zabbix.com/plugins/net/netif"
        _ "zabbix.com/plugins/net/tcp"
        ...
        _ "zabbix.com/plugins/ssl_expire"

```
### 编译安装zabbix agent2
```
yum install automake autoconf pcre* -y
./bootstrap.sh 
pushd . 
cd src/go/ 
go mod vendor 
popd 
./configure --enable-agent2 --enable-static 
make install
```
### 编辑配置文件
这里我调整了日志级别，方便前台调试
可选参数
Plugins.MQTTSubscribe.Timeout = 5
```go
egrep -v "^$|^#" conf/zabbix_agent2.conf  
LogType=console
LogFile=/tmp/zabbix_agent2.log
DebugLevel=4
Server=172.17.0.5
Plugins.Ssl_expire.Timeout=5
Hostname=node2
ControlSocket=/tmp/agent.sock
```
### **启动Zabbix_agent2**
```go
cd /root/zabbix_agent/src/go/bin
zabbix_agent2 -c conf/zabbix_agent2.conf
```
![image.png](https://cdn.nlark.com/yuque/0/2021/png/704071/1615124382567-4e083ee8-d137-428c-b5b2-0bd6ce4d511a.png#align=left&display=inline&height=81&margin=%5Bobject%20Object%5D&name=image.png&originHeight=161&originWidth=1060&size=30987&status=done&style=none&width=530)
### Zabbix创建监控项
键值示例如下
```go
ssl_expire["www.xyzabbix.cn"]
```
或
```go
ssl_expire["https://www.xyzabbix.cn"]
```
![image.png](https://cdn.nlark.com/yuque/0/2021/png/704071/1615124660164-7801e214-2bcc-497d-a252-6af930643746.png#align=left&display=inline&height=372&margin=%5Bobject%20Object%5D&name=image.png&originHeight=743&originWidth=979&size=49019&status=done&style=none&width=489.5)
查看最新数据，这个证书还有四十天过期
![image.png](https://cdn.nlark.com/yuque/0/2021/png/704071/1615124725159-142ea72a-dc40-4bfd-bd61-0af90f712ad4.png#align=left&display=inline&height=80&margin=%5Bobject%20Object%5D&name=image.png&originHeight=159&originWidth=1625&size=15012&status=done&style=none&width=812.5)
我是用的阿里云ssl证书，可以看到确实离过期时间还有四十天，今天是2021.3.7
![image.png](https://cdn.nlark.com/yuque/0/2021/png/704071/1615124786163-7b94abf7-716e-4c69-afb2-5751b053efbf.png#align=left&display=inline&height=118&margin=%5Bobject%20Object%5D&name=image.png&originHeight=236&originWidth=1497&size=21620&status=done&style=none&width=748.5)
可以创建一个触发器，在还有一个月的时候发送报警通知
![image.png](https://cdn.nlark.com/yuque/0/2021/png/704071/1615124937082-c642df97-63c0-4edc-957c-cfd8667b5d18.png#align=left&display=inline&height=367&margin=%5Bobject%20Object%5D&name=image.png&originHeight=734&originWidth=830&size=41332&status=done&style=none&width=415)
