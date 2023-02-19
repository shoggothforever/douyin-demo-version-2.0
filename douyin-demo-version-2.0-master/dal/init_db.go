package dal

import "douyin.core/config"

func InitDB() {
	config.Init()
	InitMysql()
	InitRedis()
}
