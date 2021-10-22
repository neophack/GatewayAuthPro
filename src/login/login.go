package login

import (
	"GatewayAuth/src/bindata"
	"GatewayAuth/src/config"
	"GatewayAuth/src/util"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var LoginSession = make(map[string]string)
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
	for _, v := range proxyAuth.Auth {
		if LoginSession[cookieValue] == v.Account {
			return false, nil
		}
	}
	return true, nil
}

func HttpLogin(conf config.Config) {
	//httpHandle("/static/", "http://127.0.0.1:3000")
	//httpHandle("/manifest.json", "http://127.0.0.1:3000")
	//httpHandle("/favicon.ico", "http://127.0.0.1:3000")
	//httpHandle("/logo192.png", "http://127.0.0.1:3000")
	//httpHandle("/sockjs-node", "http://127.0.0.1:3000")

	http.Handle("/login", http.FileServer(bindata.AssetFile()))

	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			s, _ := ioutil.ReadAll(r.Body) //把	body 内容读入字符串 s
			m := make(map[string]string)
			if err := json.Unmarshal(s, &m); err != nil {
				log.Println(err)
				w.Write([]byte(`{"code":0,"msg":"解析错误"}`))
				return
			}
			for _, v := range conf.Auth {
				md5str := md5Str(v.Password)
				if v.Account == m["account"] && md5str == m["password"] {
					session := md5Str(strconv.FormatInt(time.Now().UnixNano(), 10))
					LoginSession[session] = v.Account

					expiration := time.Now().Add(2 * time.Hour)
					http.SetCookie(w, &http.Cookie{Name: CookieKey, Path: "/", Value: session, Expires: expiration})

					w.Write([]byte(`{"code":200,"msg":"登录成功"}`))
					return
				}
			}
			w.Write([]byte(`{"code":0,"msg":"账号密码错误"}`))
		} else {
			w.Write([]byte("error"))
		}
	})

	http.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {

		u := util.GetUrlArg(r, "url")
		if u == "" {
			u = "/"
		}
		var param = url.Values{}
		param.Add("url", u)

		expiration := time.Now().Add(1 * time.Second)
		http.SetCookie(w, &http.Cookie{Name: CookieKey, Path: "/", Value: "", Expires: expiration})
		w.Header().Set("Location", "/login?"+param.Encode())
		w.WriteHeader(302)
	})
}

func httpHandle(pattern, targetHost string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		p, err := NewProxy(targetHost)
		if err != nil {
			panic(err)
		}
		p.ServeHTTP(w, r)
	})
}

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

func md5Str(source string) string {
	data := []byte(source)
	has := md5.Sum(data)
	return strings.ToUpper(fmt.Sprintf("%x", has))
}
