package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		_, err := fmt.Fprint(w, "<h1>Hello，这里是 goblog！</h1>")
		if err != nil {
			return
		}
	} else if r.URL.Path == "/about" {
		_, err := fmt.Fprint(w, "此博客是用以记录变成笔记，如您有反馈或建议，请联系 "+"<a href=\"mailto:zhaiyuxin103@hotmail.com\">zhaiyuxin103@hotmail.com</a>")
		if err != nil {
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, err := fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+"<p>如有疑惑，请联系我们。</p>")
		if err != nil {
			return
		}
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		return
	}
}
