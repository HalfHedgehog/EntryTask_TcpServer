package global

import (
	"TcpServer/src/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	DBHelper    *gorm.DB
	RedisHelper *redis.Client
	Config      config.Cfg
)
