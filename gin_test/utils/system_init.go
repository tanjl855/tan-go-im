package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

// InitConfig viper初始化配置信息
func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("read config err: %v\n", err)
		return
	}
	fmt.Println("config mysql:", viper.Get("mysql"))
	fmt.Println("config dsn:", viper.GetString("mysql.dsn"))
}

func InitMysql() {
	// SQL日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值
			LogLevel:      logger.Info, // 级别
			Colorful:      true,        // 彩色
		},
	)
	var err error
	DB, err = gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}
}
