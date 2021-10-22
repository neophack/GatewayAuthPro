package proxy

import (
	"GatewayAuth/src/config"
	"GatewayAuth/src/login"
	"GatewayAuth/src/util"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// NewProxy takes target host and creates a reverse proxy
// NewProxy 拿到 targetHost 后，创建一个反向代理
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

// ProxyRequestHandler handles the http request using proxy
// ProxyRequestHandler 使用 proxy 处理请求
func ProxyRequestHandler(proxyAuth config.ProxyAuth, proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		needLogin, err := login.Login(proxyAuth, r)
		if err != nil {
			log.Println(err)
		}
		if needLogin {
			u := util.GetURL(r)
			var param = url.Values{}
			param.Add("url", u)

			expiration := time.Now().Add(1 * time.Second)
			http.SetCookie(w, &http.Cookie{Name: login.CookieKey, Path: "/", Value: "", Expires: expiration})
			w.Header().Set("Location", "/login?"+param.Encode())
			w.WriteHeader(302)
		} else {
			proxy.ServeHTTP(w, r)
		}
	}
}
