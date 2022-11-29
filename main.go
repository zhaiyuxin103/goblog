package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"net/http"
	"strconv"
	"strings"
)

var router *mux.Router
var db *sql.DB

// Article 对应一条文章数据
type Article struct {
	Title, Body string
	ID          int64
}

// Delete 方法用以从数据库中山吹单条记录
func (a Article) Delete() (rowsAffected int64, err error) {
	rs, err := db.Exec("DELETE FROM articles WHERE id = " + strconv.FormatInt(a.ID, 10))

	if err != nil {
		return 0, err
	}

	// ✅删除成功
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}

	return 0, nil
}

func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

// getRouteVariable 获取 URI 路由参数
func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {

	// 1. 获取 URL 参数
	id := getRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			_, err := fmt.Fprint(w, "404 文章未找到")
			logger.LogError(err)
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprint(w, "500 服务器内部错误")
			logger.LogError(err)
		}
	} else {
		// 4. 未出现错误，执行删除操作
		rowsAffected, err := article.Delete()

		// 4.1 发生错误
		if err != nil {
			// 应该是 SQL 报错了
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprint(w, "500 服务器内部错误")
			logger.LogError(err)
		} else {
			// 4.2 未发生错误
			if rowsAffected > 0 {
				// 重定向到文章列表页
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				// Edge case
				w.WriteHeader(http.StatusNotFound)
				_, err := fmt.Fprint(w, "404 文章未找到")
				logger.LogError(err)
			}
		}
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

	database.Initialize()
	db = database.DB

	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()

	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	// 通过命名路由获取 URL 示例
	homeURL, _ := router.Get("home").URL()
	fmt.Println("homeURL：", homeURL)
	articleURL, _ := router.Get("articles.show").URL("id", "23")
	fmt.Println("articleURL：", articleURL)

	err := http.ListenAndServe(":3000", removeTrailingSlash(router))
	logger.LogError(err)
}
