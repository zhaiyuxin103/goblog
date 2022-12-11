package view

import (
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"html/template"
	"io"
	"path/filepath"
	"strings"
)

// Render 渲染视图
func Render(w io.Writer, data interface{}, tplFiles ...string) {
	// 1. 设置模版相对路径
	viewDir := "resources/views/"

	// 2. 遍历传参文件列表 Slice，设置正确的路径，支持 dir.filename 语法糖
	for i, f := range tplFiles {
		tplFiles[i] = viewDir + strings.Replace(f, ".", "/", -1) + ".tmpl"
	}

	// 3. 所有布局模版文件 Slice
	layoutFiles, err := filepath.Glob(viewDir + "layouts/*.tmpl")
	logger.LogError(err)

	// 4. 合并所有文件
	allFiles := append(layoutFiles, tplFiles...)

	// 5. 解析所有模版文件
	tmpl, err := template.New("").
		Funcs(template.FuncMap{
			"RouteName2URL": route.Name2URL,
		}).ParseFiles(allFiles...)
	logger.LogError(err)

	// 6. 渲染模版
	err = tmpl.ExecuteTemplate(w, "app", data)
	logger.LogError(err)
}
