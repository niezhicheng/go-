package initialize

import (
	"awesomeProject8/user_srv/global"
	"awesomeProject8/user_srv/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitDB()  {
	dbinfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",global.ServerConfig.MysqlInfo.User,global.ServerConfig.MysqlInfo.Password,global.ServerConfig.MysqlInfo.Host,global.ServerConfig.MysqlInfo.Port,global.ServerConfig.MysqlInfo.Name)
	fmt.Println(dbinfo,"这是沙湖")
	//dsn := "usersrv:usersrv@tcp(127.0.0.1:3306)/usersrv?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(log.New(os.Stdout,"\r\n",log.LstdFlags),logger.Config{
		SlowThreshold:             time.Second,
		Colorful:                  true,
		IgnoreRecordNotFoundError: false,
		LogLevel:                  logger.Info,
	})
	var err error
	global.DB, err = gorm.Open(mysql.Open(dbinfo), &gorm.Config{
		Logger: newLogger,
	})
	sqlDB, err := global.DB.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err != nil{
		panic(err)
	}
	_ = global.DB.AutoMigrate(&model.User{})
}
