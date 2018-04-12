sshh
[![Go Report Card](https://goreportcard.com/badge/github.com/qichengzx/sshh)](https://goreportcard.com/report/github.com/qichengzx/sshh)
----

通过使用配置文件来避免 ssh 登录服务器时需要输入密码。

Installation
------

```shell
go get github.com/qichengzx/sshh
```

Configuration
------

```json
[
  {
    "name": "test1",
    "method": "password",
    "user": "www",
    "ip": "192.168.3.10",
    "port": 22,
    "password": "yourpassword"
  },
  {
    "name": "test2",
    "method": "pem",
    "user": "root",
    "ip": "192.168.3.11",
    "port": 22,
    "password": "your pem file password or empty",
    "key": "your pem file path"
  }
]

```

Usage
------

```shell
sshh -c path-to-servers-config.json

```

配置文件路径可以是相对 sshh 的路径或绝对路径

pem 模式时，如果 pem 文件路径填写的是相对路径，程序会自动在配置文件目录和 sshh 执行目录查找，查找不到会panic。

