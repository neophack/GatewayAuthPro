package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"net/url"
)

type Config struct {
	Base  Base              `toml:"base"`
	Proxy map[string]*Proxy `toml:"proxy"`
	Auth  map[string]*Auth  `toml:"auth"`
}

type Base struct {
	Port      int      `toml:"port"`
	ProxySort []string `toml:"proxySort"`
	Crt      string       `toml:"crt"`
	Key      string       `toml:"key"`
}

type Proxy struct {
	Path        string   `toml:"path"`
	Target      string   `toml:"target"`
	CacheMaxAge int64    `toml:"cacheMaxAge"`
	HttpAuth    []string `toml:httpAuth`
	WsAuth      []string `toml:wsAuth`
}

type Auth struct {
	Account  string `toml:"account"`
	Password string `toml:"password"`
}

type ApiProxy struct {
	Path     string
	Target   *url.URL
	HttpAuth []string
	WsAuth   []string
	CacheMaxAge int64
}

func Get(filePath string) Config {
	conf := Config{
		Base: Base{
			Port: 8080,
		},
	}
	if _, err := toml.DecodeFile(filePath, &conf); err != nil {
		log.Fatalln("Config error", err)
	}
	return conf
}
