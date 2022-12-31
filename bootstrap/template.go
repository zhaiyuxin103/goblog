package bootstrap

import (
	"embed"
	"goblog/pkg/view"
)

// SetupTemplate 模版初始化
func SetupTemplate(tmplFS embed.FS) {

	view.TplFS = tmplFS

}
