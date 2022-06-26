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

	Start()

	conf := config.Get(configPath)
	jbyte, _ := json.Marshal(conf)
	log.Println(string(jbyte))

	ConfigProxy(conf)

	login.HttpLogin(conf)

	log.Println("listen : " + strconv.Itoa(conf.Base.Port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.Base.Port), nil))
}

func Start() {
	flag.StringVar(&configPath, "c", "./config", "--c config file path / 配置文件路径")
	flag.Parse()
}

func ConfigProxy(conf config.Config) {
	var x = []string{}
	var y = []string{}
	for _, i := range conf.Base.ProxySort {
		n := conf.Proxy[i]
		path := n.Path
		if len(x) == 0 {
			x = append(x, path)
			y = append(y, i)
		} else {
			for k, v := range x {
				if path == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, path)
					y = append(y, i)
				}
			}
		}
	}
	for _, p := range x {
		ProxyFunc(conf, p)
	}

}

func ProxyFunc(conf config.Config, path string) {

	// handle all requests to your server using the proxy
	// 使用 proxy 处理所有请求到你的服务
	http.HandleFunc(path, proxy.ProxyRequestHandler(conf))
}
