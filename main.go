package main

import (
	"fmt"
	"net/http"
	"strings"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		_, err := fmt.Fprint(w, "<h1>Hello，欢迎来到 goblog！</h1>")
		if err != nil {
			return
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		_, err := fmt.Fprintf(w, "<h1>请求页面未找到 :(</h1>"+"<p>如有疑惑，请联系我们。</p>")
		if err != nil {
			return
		}
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := fmt.Fprint(w, "此博客是用以记录变成笔记，如您有反馈或建议，请联系 "+"<a href=\"mailto:zhaiyuxin103@hotmail.com\">zhaiyuxin103@hotmail.com</a>")
	if err != nil {
		return
	}
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/about", aboutHandler)

	// 文章详情
	router.HandleFunc("/articles/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.SplitN(r.URL.Path, "/", 3)[2]
		_, err := fmt.Fprint(w, "文章 ID："+id)
		if err != nil {
			return
		}
	})

	// 列表 or 创建
	router.HandleFunc("/articles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			_, err := fmt.Fprint(w, "访问文章列表")
			if err != nil {
				return
			}
		case "POST":
			_, err := fmt.Fprint(w, "创建新的文章")
			if err != nil {
				return
			}
		}
	})

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		return
	}
}
