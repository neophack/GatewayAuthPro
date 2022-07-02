package main

import (
	"GatewayAuth/src/config"
	"GatewayAuth/src/login"
	"GatewayAuth/src/proxy"
	"encoding/json"
	"flag"
	"github.com/rs/cors"
	"log"
	"net/http"
	"time"
)
//go:generate go install -a -v github.com/go-bindata/go-bindata/...@latest
//go:generate /home/nn/GoCode/bin/go-bindata -fs -o=../bindata/bindata.go -pkg=bindata ../../../frontend/build

var configPath string

func main() {

	Start()

	conf := config.Get(configPath)
	jbyte, _ := json.Marshal(conf)
	log.Println(string(jbyte))

	mux := http.NewServeMux()
	login.HttpLogin(conf,mux)

	gateway := proxy.NewApiGateway(conf)

	mux.HandleFunc("/", gateway.Handle)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Accept", "Content-Type", "Authorization"},
	}).Handler(mux)

	log.Print("gateway server started")

	srv := &http.Server{
		Addr:         ":10000",
		Handler:      handler,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func Start() {
	flag.StringVar(&configPath, "c", "./config", "--c config file path / 配置文件路径")
	flag.Parse()
}
