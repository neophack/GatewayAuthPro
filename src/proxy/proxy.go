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

type ApiGateway struct {
	Proxyinfos []config.ApiProxy
	Conf       config.Config
}

func NewApiGateway(conf config.Config) *ApiGateway {
	var proxyInfos []config.ApiProxy
	for _, p := range conf.Base.ProxySort {
		n := conf.Proxy[p]
		backendUri, err := url.Parse(n.Target)

		if err != nil || backendUri.Host == "" {
			panic("invalid BACKEND_URI setting")
		}
		var proxyInfo config.ApiProxy
		proxyInfo.HttpAuth = n.HttpAuth
		proxyInfo.WsAuth = n.WsAuth
		proxyInfo.Target = backendUri
		proxyInfo.Path = n.Path
		proxyInfo.CacheMaxAge = 0
		proxyInfos = append(proxyInfos, proxyInfo)
	}
	return &ApiGateway{
		Proxyinfos: proxyInfos,
		Conf:       conf,
	}
}

func (gateway *ApiGateway) Handle(w http.ResponseWriter, r *http.Request) {
	loginState, cacheMaxAge, target, err := login.Login(gateway.Conf, gateway.Proxyinfos, r)
	if err != nil && err.Error() != "http: named cookie not present" {
		log.Println(err)
	}
	if target == nil {
		loginState = 1
	}

	switch loginState {
	case login.NotLogin:
		login.ClearCookie(w)
		http.SetCookie(w, &http.Cookie{Name: login.CookieKey, Path: "/", Value: "", HttpOnly: true, MaxAge: -1})
		w.Header().Set("Location", "/loginxx?"+ParamEncode(r))
		w.WriteHeader(http.StatusFound)
	case login.NoPermission:
		w.Header().Set("Location", "/loginxx?type=nopermission&"+ParamEncode(r))
		w.WriteHeader(http.StatusFound)
	case login.AlreadyLogin:
		proxy := httputil.NewSingleHostReverseProxy(target)
		ServeHttp(proxy, cacheMaxAge, w, r, target)
		//proxy.ServeHTTP(w, r)
	case login.NoLogin:
		proxy := httputil.NewSingleHostReverseProxy(target)
		ServeHttp(proxy, cacheMaxAge, w, r, target)
	}

}

func ServeHttp(proxy *httputil.ReverseProxy, cacheMaxAge int64, w http.ResponseWriter, r *http.Request, u *url.URL) {
	if cacheMaxAge > 0 {
		w.Header().Set("Cache-Control", "max-age="+strconv.FormatInt(cacheMaxAge, 10))
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}

	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = u.Host

	proxy.ServeHTTP(w, r)
}

func ParamEncode(r *http.Request) string {

	u := util.GetURL(r)
	var param = url.Values{}
	param.Add("url", u)

	return param.Encode()
}
