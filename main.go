package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

var router = mux.NewRouter()

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "<h1>Hello，欢迎来到 goblog！</h1>")
	if err != nil {
		return
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+"<a href=\"mailto:zhaiyuxin103@hotmail.com\">zhaiyuxin103@hotmail.com</a>")
	if err != nil {
		return
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
	if err != nil {
		return
	}
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "访问文章列表")
	if err != nil {
		return
	}
}

func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>创建文章 -- 我的技术博客</title>
</head>
<body>
	<form action="%s" method="post">
		<p><input type="text" name="title"></p>
		<p><textarea name="body" cols="30" rows="10"></textarea></p>
		<p><button type="submit">提交</button></p>
	</form>
</body>
</html>
`
	storeURL, _ := router.Get("articles.store").URL()
	fmt.Fprintf(w, html, storeURL)
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "创建新的文章")
	if err != nil {
		return
	}
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := fmt.Fprint(w, "文章 ID："+id)
	if err != nil {
		return
	}
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 2. 继续处理请求
		next.ServeHTTP(w, r)
	})
}

func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 除首页以外，移除所有请求路径后面的斜杠
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		// 2. 将请求传递下去
		next.ServeHTTP(w, r)
	})
}

func main() {
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("articles.create")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")

	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	// 通过命名路由获取 URL 示例
	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL：", homeURL)
	articleURL, _ := router.Get("articles.show").URL("id", "23")
	fmt.Println("articleURL：", articleURL)

	err := http.ListenAndServe(":3000", removeTrailingSlash(router))
	if err != nil {
		return
	}
}
