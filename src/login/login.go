package login

import (
	"GatewayAuth/src/config"
	"net/http"
	"strings"
)

var LoginSession map[string]string
var CookieKey = "gateway.auto.cookie"

func Login(proxyAuth config.ProxyAuth, r *http.Request) (needLogin bool, err error) {
	upgrade := r.Header.Get("Upgrade")
	if upgrade == "" {
		upgrade = r.Header.Get("upgrade")
	}
	if upgrade == "" {
		upgrade = r.Header.Get("UPGRADE")
	}
	var auth []string
	if upgrade == "websocket" || upgrade == "Websocket" || upgrade == "WEBSOCKET" {
		auth = proxyAuth.WsAuth
	} else {
		auth = proxyAuth.HttpAuth
	}
	if len(auth) <= 0 {
		return false, nil
	}
	var cookie *http.Cookie
	if cookie, err = r.Cookie(CookieKey); err != nil {
		return true, err
	}
	cookieValue := strings.TrimSpace(cookie.Value)
	if cookieValue == "" {
		return true, nil
	}
	if LoginSession[cookieValue] == proxyAuth.Hash {
		return false, nil
	} else {
		return true, nil
	}
}

func HttpLogin() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Write([]byte( "success"))
		} else if r.Method == "POST" {

		} else {
			w.Write([]byte( "error"))
		}
	})
}
