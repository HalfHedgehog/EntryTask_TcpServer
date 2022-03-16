package initialize

import (
	"TcpServer/src/config"
	"TcpServer/src/global"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitConfig 初始化配置文件
func InitConfig() {
	viper.AddConfigPath("./src")
	viper.SetConfigName("apps")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("read config file failed, %v", err)
	}
	c := config.Cfg{}
	if err := viper.Unmarshal(&c); err != nil {
		fmt.Printf("unmarshal config file failed, %v", err)
	}
	global.Config = c
}

// CreateGorm 初始化Gorm
func CreateGorm() (db *gorm.DB) {
	user := global.Config.Mysql.User
	password := global.Config.Mysql.Password
	address := global.Config.Mysql.Address
	port := global.Config.Mysql.Port
	dbName := global.Config.Mysql.DB

	dsn :=
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, address, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func CreateRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     global.Config.Redis.Address,
		Password: global.Config.Redis.Password,
		DB:       global.Config.Redis.DB,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return client
}
