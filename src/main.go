package main

import (
	"GatewayAuth/src/config"
	"GatewayAuth/src/login"
	"GatewayAuth/src/proxy"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func main() {

	conf := config.Get("./config")
	jbyte, _ := json.Marshal(conf)
	log.Println(string(jbyte))

	configProxy(conf)

	login.HttpLogin()

	log.Println("listen : " + strconv.Itoa(conf.Base.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.Base.Port), nil))
}

func configProxy(conf config.Config) {

	for p := range conf.Proxy {
		n := conf.Proxy[p]
		proxyAuth := config.ToProxyAuth(*n, conf.Auth)
		for e := range n.Path {
			path := n.Path[e]
			for e2 := range n.Target {
				target := n.Target[e2]
				proxyFunc(proxyAuth, path, target)
			}
		}
	}

}

func proxyFunc(proxyAuth config.ProxyAuth, path, target string) {
	proxy2, err := proxy.NewProxy(target)
	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	// 使用 proxy 处理所有请求到你的服务
	http.HandleFunc(path, proxy.ProxyRequestHandler(proxyAuth, proxy2))
}
