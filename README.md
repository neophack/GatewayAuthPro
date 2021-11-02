<p align="center">
  <a href="https://github.com/LambdaExpression/GatewayAuth">
    <img width="150" src="/public/logo.png">
  </a>
</p>

<h1 align="center">Gateway Auth</h1>

<div align="center">Go / React / Material-ui 网关登录工具 </div>
<br/>
<div align="center">
    <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/LambdaExpression/GatewayAuth?cacheSeconds=3600">
    <a href="javascript:void(0)"><img src="https://img.shields.io/badge/Node-v14.15.4-blue.svg" /></a>
    <a href="javascript:void(0)"><img src="https://img.shields.io/badge/React-17.0.2-blue.svg" /></a>
    <a href="javascript:void(0)"><img src="https://img.shields.io/badge/Material-5.0.4-blue.svg" /></a>
    <!--<img alt="GitHub repo file count" src="https://img.shields.io/github/directory-file-count/LambdaExpression/GatewayAuth">-->
    <img alt="GitHub repo size" src="https://img.shields.io/github/repo-size/LambdaExpression/GatewayAuth">
    <a href="https://codeclimate.com/github/LambdaExpression/GatewayAuth/maintainability"><img src="https://api.codeclimate.com/v1/badges/027ea75f1d5cbf0a3081/maintainability" /></a>
</div>
<br/>

<br/>
<br/>

<p align="center">
  <img width="49%" src="/public/image1.png">
  <img width="49%" src="/public/image2.png">
</p>




### Run

[Download the version of the corresponding system](https://github.com/LambdaExpression/GatewayAuth/releases)

```sh
$ wget https://github.com/LambdaExpression/GatewayAuth/releases/download/v1.0.0/gatewayAuth_linux_amd64_1_0_0

$ chmod +x gatewayAuth_linux_amd64_1_0_0

$ ./gatewayAuth_darwin_amd64_1_0_0 -h
Usage of ./gatewayAuth_darwin_amd64_1_0_0:
  -c string
    	--c config file path / 配置文件路径 (default "./config")
    	
$ echo -e '[base]\nport = 8094\nproxySort=["test"]\n[proxy]\n    [proxy.test]\n    path = "/"\n    target = "http://127.0.0.1:80"\n    httpAuth = ["tom"]\n[auth]\n    [auth.test]\n    account = "test"\n    password = "123"' > ./config

$ ./gatewayAuth_darwin_amd64_1_0_0 -c ./config
2021/11/01 16:13:16 {"Base":{"Port":8094,"ProxySort":["test","serverstatusws","serverstatus"]},"Proxy":{"serverstatus":{"Path":"/","Target":"http://127.0.0.1:35601","CacheMaxAge":0,"HttpAuth":["tom"],"WsAuth":null},"serverstatusws":{"Path":"/public","Target":"http://127.0.0.1:35601","CacheMaxAge":0,"HttpAuth":null,"WsAuth":["tom"]},"test":{"Path":"/test","Target":"http://127.0.0.1:80","CacheMaxAge":0,"HttpAuth":["tom"],"WsAuth":["tom"]}},"Auth":{"test":{"Account":"test","Password":"123"},"tom":{"Account":"tom","Password":"123"}}}
2021/11/01 16:13:16 listen : 8094

```


### Config File
```toml
[base]     
port = 8094
# proxy execution order / 代理执行顺序
proxySort=["test","serverstatusws","serverstatus"] 

[proxy]

    [proxy.test]
    path = "/test"
    target = "http://127.0.0.1:80"
    httpAuth = ["tom"]   # login account / 登录账号
    wsAuth = ["tom"]     # login account / 登录账号

    [proxy.serverstatusws]
    path = "/public"
    target = "http://127.0.0.1:35601"
    wsAuth = ["tom"]

    [proxy.serverstatus]
    path = "/"
    target = "http://127.0.0.1:35601"
    httpAuth = ["tom"]

[auth]
    
    # account and password / 账号密码
    [auth.tom]
    account = "tom"
    password = "123"

    [auth.test]
    account = "test"
    password = "123"
```

