package router

import (
	"awesomeProject8/mxshop-api/user-web/api"
	"github.com/gin-gonic/gin"
)

func InitBaseRouter(group *gin.RouterGroup) {
	BaseRouter := group.Group("base")
	{
		BaseRouter.GET("captcha",api.Captcha)
		BaseRouter.POST("send_sms",api.SendSmg)
		BaseRouter.POST("register",api.Register)
	}
}
