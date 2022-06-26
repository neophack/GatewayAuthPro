package proxy

import (
	"GatewayAuth/src/config"
	"GatewayAuth/src/login"
	"GatewayAuth/src/util"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

// NewProxy takes target host and creates a reverse proxy
// NewProxy 拿到 targetHost 后，创建一个反向代理
func NewProxy(targetHost string) (*url.URL, error) {
	urlTarget, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return urlTarget, nil
}

// ProxyRequestHandler handles the http request using proxy
// ProxyRequestHandler 使用 proxy 处理请求
func ProxyRequestHandler(conf config.Config) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		loginState, cacheMaxAge, target, err := login.Login(conf, r)
		if err != nil && err.Error() != "http: named cookie not present" {
			log.Println(err)
		}
		urlTarget, err := NewProxy(target)
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(urlTarget)

		switch loginState {
		case login.NotLogin:
			login.ClearCookie(w)
			http.SetCookie(w, &http.Cookie{Name: login.CookieKey, Path: "/", Value: "", HttpOnly: true, MaxAge: -1})
			w.Header().Set("Location", "/login?"+ParamEncode(r))
			w.WriteHeader(http.StatusFound)
		case login.NoPermission:
			w.Header().Set("Location", "/login?type=nopermission&"+ParamEncode(r))
			w.WriteHeader(http.StatusFound)
		case login.AlreadyLogin:
			ServeHttp(proxy, cacheMaxAge, w, r, urlTarget)
		case login.NoLogin:
			ServeHttp(proxy, cacheMaxAge, w, r, urlTarget)
		}
	}
}

func ServeHttp(proxy *httputil.ReverseProxy, cacheMaxAge int64, w http.ResponseWriter, r *http.Request, u *url.URL) {
	if cacheMaxAge > 0 {
		w.Header().Set("Cache-Control", "max-age="+strconv.FormatInt(cacheMaxAge, 10))
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}

	w.Header().Set("Sec-WebSocket-Accept", r.Header.Get("Sec-WebSocket-Accept"))

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = u.Host
	w.Header().Set("Access-Control-Allow-Origin", "*")
	proxy.ServeHTTP(w, r)
}

func ParamEncode(r *http.Request) string {

	u := util.GetURL(r)
	var param = url.Values{}
	param.Add("url", u)

	return param.Encode()
}
