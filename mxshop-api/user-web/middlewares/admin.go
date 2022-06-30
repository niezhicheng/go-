package middlewares

import (
	"awesomeProject8/mxshop-api/user-web/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsadminAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		cliasa,_ := context.Get("claims")
		currentUser := cliasa.(*models.CustomClaims).AuthorityId
		fmt.Println(currentUser)
		if currentUser == 1 {
			context.JSON(http.StatusForbidden,gin.H{
				"msg": "无权限",
			})
			context.Abort()
			return
		}
		context.Next()
	}
}
