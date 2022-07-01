package global

import (
	"awesomeProject8/goods_srv/config"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	ServerConfig *config.ServerConfig
)

