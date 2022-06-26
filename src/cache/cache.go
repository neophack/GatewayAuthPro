package cache

import "time"

var cacheSize = 100
var cacheMaps [100]map[string]string

func init() {
	for i := 0; i < cacheSize; i++ {
		cacheMaps[i] = make(map[string]string)
	}
}

func Get(key string) string {
	return Find(func(i int) string {
		return cacheMaps[i][key]
	})
}

func Set(key, value string) {
	Find(func(i int) string {
		cacheMaps[i][key] = value
		if (i + 1) >= len(cacheMaps) {
			cacheMaps[0][key] = value
		} else {
			cacheMaps[i+1][key] = value
		}
		return ""
	})
}

func Find(f func(i int) string) string {
	t := GetTime()
	var v string
	for i := 0; i < cacheSize; i++ {
		if t == i {
			v = f(i)
		} else if (t + 1) == i {
		} else {
			cacheMaps[i] = make(map[string]string)
		}
	}
	return v
}

func GetTime() int {
	u := time.Now().Unix()
	u = u / 60 / 60
	return int(u % int64(cacheSize))
}
