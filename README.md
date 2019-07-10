sshh
[![Go Report Card](https://goreportcard.com/badge/github.com/qichengzx/sshh)](https://goreportcard.com/report/github.com/qichengzx/sshh)
----

用Go写的一个通过读取配置文件连接远程服务器的工具，免去ssh登录时需要输入密码的步骤。

Installation
------

```shell
go get github.com/qichengzx/sshh
```

or Download

[https://github.com/qichengzx/sshh/releases](https://github.com/qichengzx/sshh/releases)

Build for Multi Platforms
------

```shell
chmod +x build.sh
sh build.sh
```

Build for current Platform
------
```shell
go build main.go
```

将会在releases目录生成可执行文件

Configuration
------

```yaml
group:
  -
    -
      name: g1
      method: password
      user: www
      port: 22
      password: yourpassword
      ip: [192.168.3.10, 192.168.3.11, 192.168.3.12]
      key:
single:
  -
    name: test1
    method: password
    user: www
    ip: 192.168.3.20
    port: 22
    password: yourpassword
  -
    name: test2
    method: pem
    user: root
    ip: 192.168.3.21
    port: 22
    password: 'your pem file password or empty'
    key: 'your pem file path'
```

其中，group 表示可以使用一个账号登录的一组服务器，single 则为每服务器独立账号或密码的情况。
(You know,Sometimes you wake up and you have a set of server permissions.)

Usage
------

第一次使用需要指定配置文件：

```shell
sshh -c path-to-servers-config.yaml
```

之后如果配置不变的情况下，只需要 ```sshh```  即可。

配置文件路径可以是相对 sshh 的路径或绝对路径

pem 模式时，如果 pem 文件路径填写的是相对路径，程序会自动在配置文件目录和 sshh 执行目录查找，查找不到会panic。