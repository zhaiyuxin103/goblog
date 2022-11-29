package controllers

import (
	"fmt"
	"goblog/pkg/logger"
	"net/http"
)

// PagesController 处理静态页面
type PagesController struct {
}

// Home 首页
func (*PagesController) Home(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "<h1>Hello，欢迎来到 goblog！</h1>")
	logger.LogError(err)
}

// About 关于我们页面
func (*PagesController) About(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+"<a href=\"mailto:zhaiyuxin103@hotmail.com\">zhaiyuxin103@hotmail.com</a>")
	logger.LogError(err)
}

// NotFound 404 页面
func (*PagesController) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, err := fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
	logger.LogError(err)
}
