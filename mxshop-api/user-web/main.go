package main

import (
	"awesomeProject8/mxshop-api/user-web/global"
	"awesomeProject8/mxshop-api/user-web/initalize"
	"awesomeProject8/mxshop-api/user-web/utils"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	validatatorcur "awesomeProject8/mxshop-api/user-web/validatator"
)


func main()  {
	initalize.InitLogeer()
	initalize.InitConfig()
	fmt.Println(global.ServerConfig.UserSrvInfo.Port,"这是")
	fmt.Println(global.ServerConfig.JWTInfo.SigningKey)
	Routers := initalize.Routers()
	if err := initalize.InitTrans("zh"); err != nil{
		zap.S().Error(err.Error())
	}
	if v,ok := binding.Validator.Engine().(*validator.Validate); ok{
		_ = v.RegisterValidation("mobile", validatatorcur.ValidateMobile)
		_ = v.RegisterTranslation("mobile",global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile","{0} 非法的手机号码!",true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t,_ := ut.T("mobile",fe.Field())
			return t
		})
	}
	zap.S().Infof("启动服务器 端口: %d",global.ServerConfig.Port)
	port,err := utils.GetFreePort()
	if err != nil{
		zap.S().Errorw("错误了")
	}
	global.ServerConfig.Port = port
	if err := Routers.Run(fmt.Sprintf(":%d",global.ServerConfig.Port));err != nil{
		zap.S().Panic("启动失败",err.Error())
	}
}
