package controllers

import (
	"fmt"
	"goblog/app/models/article"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"unicode/utf8"
)

// ArticlesController 文章相关页面
type ArticlesController struct {
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
		// ---  2. 加载模版 ---
		view.Render(w, view.D{
			"Articles": articles,
		}, "articles.index")
	}
}

// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "articles.create", "articles._form_field")
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
		view.Render(w, view.D{
			"Title":  title,
			"Body":   body,
			"Errors": errors,
		}, "articles.create", "articles._form_field")
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
		// ---  4. 读取成功，显示文章 ---
		view.Render(w, view.D{
			"Article": article,
		}, "articles.show")
	}
}

// Edit 文章更新页面
func (*ArticlesController) Edit(w http.ResponseWriter, r *http.Request) {

	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

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
		// 4. 读取成功，显示编辑文章表单页
		view.Render(w, view.D{
			"Title":   _article.Title,
			"Body":    _article.Body,
			"Article": _article,
			"Errors":  nil,
		}, "articles.edit", "articles._form_field")
	}
}

// Update 更新文章
func (*ArticlesController) Update(w http.ResponseWriter, r *http.Request) {

	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

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
		// 4. 未出现错误

		// 4.1 表单验证
		title := r.PostFormValue("title")
		body := r.PostFormValue("body")

		errors := validateArticleFormData(title, body)

		if len(errors) == 0 {

			// 4.2 表单验证通过，更新数据
			_article.Title = title
			_article.Body = body

			rowsAffected, err := _article.Update()

			if err != nil {
				logger.LogError(err)
				w.WriteHeader(http.StatusInternalServerError)
				_, err := fmt.Fprint(w, "500 服务器内部错误")
				logger.LogError(err)
			}

			// ✅更新成功，跳转到文章详情页
			if rowsAffected > 0 {
				showURL := route.Name2URL("articles.show", "id", id)
				http.Redirect(w, r, showURL, http.StatusFound)
			} else {
				_, err := fmt.Fprint(w, "您没有做任何更改！")
				logger.LogError(err)
			}
		} else {

			// 4.3 表单验证不通过，显示理由
			view.Render(w, view.D{
				"Title":   title,
				"Body":    body,
				"Article": _article,
				"Errors":  errors,
			}, "articles.edit", "articles._form_field")
		}
	}
}

// Delete 删除文章
func (*ArticlesController) Delete(w http.ResponseWriter, r *http.Request) {

	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", r)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

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
		// 4. 未出现错误，执行删除操作
		rowsAffected, err := _article.Delete()

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
				indexURL := route.Name2URL("articles.index")
				http.Redirect(w, r, indexURL, http.StatusFound)
			} else {
				// Edge case
				w.WriteHeader(http.StatusNotFound)
				_, err := fmt.Fprint(w, "404 文章未找到")
				logger.LogError(err)
			}
		}
	}
}
