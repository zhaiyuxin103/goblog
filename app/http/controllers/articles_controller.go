package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/types"
	"gorm.io/gorm"
	"html/template"
	"net/http"
	"strconv"
	"unicode/utf8"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
}

// ArticlesFormData 创建博文表单数据
type ArticlesFormData struct {
	Title, Body string
	URL         string
	Errors      map[string]string
}

func validateArticleFormData(title string, body string) map[string]string {
	errors := make(map[string]string)
	// 验证标题
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 40 {
		errors["title"] = "标题长度需介于 3 ～ 40 之间"
	}

	// 验证内容
	if body == "" {
		errors["body"] = "内容不能为空"
	} else if utf8.RuneCountInString(body) < 10 {
		errors["body"] = "内容长度需大于或等于 10 个字节"
	}

	return errors
}

// Index 文章列表页
func (*ArticlesController) Index(w http.ResponseWriter, r *http.Request) {
	// 1. 获取结果集
	articles, err := article.GetAll()

	if err != nil {
		// 数据库错误
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprint(w, "500 服务器内部错误")
		logger.LogError(err)
	} else {
		// 2. 加载模版
		tmpl, err := template.ParseFiles("resources/views/articles/index.tmpl")
		logger.LogError(err)

		// 3. 渲染模版，将所有文章的数据传输进去
		err = tmpl.Execute(w, articles)
		logger.LogError(err)
	}
}

// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {

	storeURL := route.Name2URL("articles.store")
	data := ArticlesFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}
	tmpl, err := template.ParseFiles("resources/views/articles/create.tmpl")
	logger.LogError(err)

	err = tmpl.Execute(w, data)
	logger.LogError(err)
}

// Store 创建文章
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	title := r.PostFormValue("title")
	body := r.PostFormValue("body")

	errors := validateArticleFormData(title, body)

	// 检查是否有错误
	if len(errors) == 0 {
		_article := article.Article{
			Title: title,
			Body:  body,
		}
		err := _article.Create()
		if _article.ID > 0 {
			_, err := fmt.Fprint(w, "插入成功，ID 为："+strconv.FormatUint(_article.ID, 10))
			logger.LogError(err)
		} else {
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err := fmt.Fprint(w, "500 服务器内部错误")
			logger.LogError(err)
		}
	} else {
		storeURL := route.Name2URL("articles.store")

		data := ArticlesFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}
		tmpl, err := template.ParseFiles("resources/views/articles/create.tmpl")
		logger.LogError(err)

		err = tmpl.Execute(w, data)
		logger.LogError(err)
	}
}

// Show 文章详情页面
func (*ArticlesController) Show(w http.ResponseWriter, r *http.Request) {
	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
		// 4. 读取成功
		tmpl, err := template.New("show.tmpl").
			Funcs(template.FuncMap{
				"RouteName2URL":  route.Name2URL,
				"Uint64ToString": types.Uint64ToString,
			}).
			ParseFiles("resources/views/articles/show.tmpl")
		logger.LogError(err)

		err = tmpl.Execute(w, article)
		logger.LogError(err)
	}
}
