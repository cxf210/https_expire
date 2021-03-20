# https_expire
zabbix-agent2自定义插件监控https证书过期时间
> 有问题可以联系微信：wanger5354

> 微信公众号：运维开发故事

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
Plugins.Https_expire.Timeout = 5
```go
egrep -v "^$|^#" conf/zabbix_agent2.conf  
LogType=console
LogFile=/tmp/zabbix_agent2.log
DebugLevel=4
Server=172.17.0.5
Plugins.Https_expire.Timeout=5
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
https_expire["www.xyzabbix.cn"]
```
或
```go
https_expire["https://www.xyzabbix.cn"]
```
![image.png](https://cdn.nlark.com/yuque/0/2021/png/704071/1615291973645-d92ab392-4c59-4739-a063-8379236b9b36.png#align=left&display=inline&height=336&margin=%5Bobject%20Object%5D&name=image.png&originHeight=672&originWidth=930&size=46093&status=done&style=none&width=465)
查看最新数据，这个证书还有四十天过期
![image.png](https://cdn.nlark.com/yuque/0/2021/png/704071/1615124725159-142ea72a-dc40-4bfd-bd61-0af90f712ad4.png#align=left&display=inline&height=80&margin=%5Bobject%20Object%5D&name=image.png&originHeight=159&originWidth=1625&size=15012&status=done&style=none&width=812.5)
我是用的阿里云ssl证书，可以看到确实离过期时间还有四十天，今天是2021.3.7
![image.png](https://cdn.nlark.com/yuque/0/2021/png/704071/1615124786163-7b94abf7-716e-4c69-afb2-5751b053efbf.png#align=left&display=inline&height=118&margin=%5Bobject%20Object%5D&name=image.png&originHeight=236&originWidth=1497&size=21620&status=done&style=none&width=748.5)
可以创建一个触发器，在还有一个月的时候发送报警通知
![image.png](https://cdn.nlark.com/yuque/0/2021/png/704071/1615292050620-890747a2-b854-4e91-b18c-c5ca272575f2.png#align=left&display=inline&height=302&margin=%5Bobject%20Object%5D&name=image.png&originHeight=603&originWidth=949&size=40242&status=done&style=none&width=474.5)
