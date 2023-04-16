package apiserver

import "github.com/gin-gonic/gin"

func initRouter(g *gin.Engine) {
	installMiddlewares(g) // 安装需要的中间件
	installController(g)  // 安装控制器
}

func installMiddlewares(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {
	// TODO: 执行安装中间件的操作

	return g
}
