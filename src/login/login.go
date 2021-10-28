package login

import (
	"GatewayAuth/src/bindata"
	"GatewayAuth/src/cache"
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

var CookieKey = "gateway.auto.cookie"

var NotLogin = 1     // 未登录
var NoPermission = 2 // 无权限
var AlreadyLogin = 3 // 已登录
var NoLogin = 4      // 免登录

func Login(conf config.Config, r *http.Request) (loginState int, cacheMaxAge int64, err error) {
	upgrade := r.Header.Get("Upgrade")
	if upgrade == "" {
		upgrade = r.Header.Get("upgrade")
	}
	if upgrade == "" {
		upgrade = r.Header.Get("UPGRADE")
	}
	isWs := upgrade == "websocket" || upgrade == "Websocket" || upgrade == "WEBSOCKET"

	for _, s := range conf.Base.ProxySort {
		v := conf.Proxy[s]
		if strings.HasPrefix(r.URL.Path, v.Path) {
			var s []string
			if isWs {
				s = v.WsAuth
				log.Println("WsAuth ", v.WsAuth)
			} else {
				s = v.HttpAuth
				log.Println("HttpAuth ", v.HttpAuth)
			}
			log.Println("isWs ", isWs)
			log.Println("v ", v)
			log.Println("s ", s)
			log.Println("r.URL.Path ", r.URL.Path)
			if len(s) == 0 {
				log.Println(1)
				return NoLogin, v.CacheMaxAge, nil
			}

			var cookie *http.Cookie
			if cookie, err = r.Cookie(CookieKey); err != nil {
				log.Println(2)
				return NotLogin, v.CacheMaxAge, err
			}
			cookieValue := strings.TrimSpace(cookie.Value)
			if cookieValue == "" {
				log.Println(3)
				return NotLogin, v.CacheMaxAge, nil
			}

			cv := cache.Get(cookieValue)
			if cv == "" {
				log.Println(4)
				return NotLogin, v.CacheMaxAge, nil
			}

			for _, v2 := range s {
				p := conf.Auth[v2]
				if cv == p.Account {
					log.Println(5)
					return AlreadyLogin, v.CacheMaxAge, nil
				}
			}
			log.Println(6)
			return NoPermission, v.CacheMaxAge, nil
		}
	}

	return NotLogin, 0, nil
}

func HttpLogin(conf config.Config) {

	http.Handle("/login/", http.StripPrefix("/login/", http.FileServer(bindata.AssetFile())))

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
					cache.Set(session, v.Account)

					expiration := time.Now().Add(2 * time.Hour)
					http.SetCookie(w, &http.Cookie{Name: CookieKey, Path: "/", Value: session, HttpOnly: true, Expires: expiration})

					w.Write([]byte(`{"code":200,"msg":"登录成功"}`))
					return
				}
			}
			w.Write([]byte(`{"code":0,"msg":"账号密码错误"}`))
		} else {
			w.Write([]byte("error"))
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {

		u := util.GetUrlArg(r, "url")
		if strings.TrimSpace(u) == "" || strings.TrimSpace(u) == "null" {
			u = "/"
		}
		var param = url.Values{}
		param.Add("url", u)
		http.SetCookie(w, &http.Cookie{Name: CookieKey, Path: "/", Value: "", HttpOnly: true, MaxAge: -1})
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

func ClearCookie(w http.ResponseWriter) {
	expiration := time.Now().Add(1 * time.Second)
	http.SetCookie(w, &http.Cookie{Name: CookieKey, Path: "/", Value: "", Expires: expiration})
}
