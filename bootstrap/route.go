package bootstrap

import (
	"github.com/gorilla/mux"
	"goblog/routes"
)

// SetupRoute 路由初始化
func SetupRoute() *mux.Router {
	router := mux.NewRouter()
	routes.RegisterWebRoutes(router)
	routes.RegisterAPIRoutes(router)
	return router
}
