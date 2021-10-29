<p align="center">
  <a href="https://github.com/LambdaExpression/GatewayAuth">
    <img width="150" src="/public/logo.png">
  </a>
</p>

<h1 align="center">Gateway Auth</h1>

<div align="center">Go/React/Material-ui 网关登录工具</div>
<br/>
<br/>

<p align="center">
  <img width="49%" src="/public/image1.png">
  <img width="49%" src="/public/image2.png">
</p>



#### config file / config 文件
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
    
    # account password / 账号密码
    [auth.tom]
    account = "tom"
    password = "123"

    [auth.test]
    account = "test"
    password = "123"
```