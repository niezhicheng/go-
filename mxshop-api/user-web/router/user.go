package router

import (
	"awesomeProject8/mxshop-api/user-web/api"
	"awesomeProject8/mxshop-api/user-web/middlewares"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitUserRouter(group *gin.RouterGroup)  {
	//UserRouter := group.Group("user")
	UserRouter := group.Group("user").Use(middlewares.JWTAuth())
	zap.S().Infof("启动了router")
	{
		//UserRouter.GET("list",api.GetUserList)
		UserRouter.GET("list",middlewares.IsadminAuth(),api.GetUserList)
		UserRouter.POST("pwdlogin",api.PasswordLogin)
		UserRouter.GET("ceshi",api.Ceshiciew)
	}
}