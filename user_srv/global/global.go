package global

import (
	"awesomeProject8/user_srv/config"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	ServerConfig *config.ServerConfig
)

