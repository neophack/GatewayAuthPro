package main

import (
	"GatewayAuth/src/config"
	"GatewayAuth/src/login"
	"GatewayAuth/src/proxy"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
)

var configPath string

func main() {

	start()

	conf := config.Get(configPath)
	jbyte, _ := json.Marshal(conf)
	log.Println(string(jbyte))

	configProxy(conf)

	login.HttpLogin(conf)

	log.Println("listen : " + strconv.Itoa(conf.Base.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.Base.Port), nil))
}

func start() {
	flag.StringVar(&configPath, "c", "./config", "--c config file path / 配置文件路径")
	flag.Parse()
}

func configProxy(conf config.Config) {

	for _, p := range conf.Base.ProxySort {
		n := conf.Proxy[p]
		path := n.Path
		proxyFunc(conf, path, n.Target)
	}

}

func proxyFunc(conf config.Config, path, target string) {
	proxy2, err := proxy.NewProxy(target)
	if err != nil {
		panic(err)
	}

	// handle all requests to your server using the proxy
	// 使用 proxy 处理所有请求到你的服务
	http.HandleFunc(path, proxy.ProxyRequestHandler(conf, proxy2))
}
