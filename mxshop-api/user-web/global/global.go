package global

import (
	"awesomeProject8/mxshop-api/user-web/config"
	"awesomeProject8/user_srv/proto"
	ut "github.com/go-playground/universal-translator"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	Trans ut.Translator
	UserSrvClient proto.UserClient
	NacosConfig *config.NacOsCOnfig = & config.NacOsCOnfig{}
)
