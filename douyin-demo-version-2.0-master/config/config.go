package config

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

const (
	Sqldsn        string = "dsn"
	Redisaddr     string = "addr"
	Redispassword string = "password"
)

type Server struct {
	IP   string
	Port int
}
type mMysql struct {
	Host      string
	Port      int
	Database  string
	Username  string
	Password  string
	Charset   string
	ParseTime bool `toml:"parse_time"`
	Loc       string
}
type mRedis struct {
	redis.Options
}
type Config struct {
	DB     mMysql `toml:"mysqldal"`
	Server `toml:"server"`
}

var Info Config

// DBConnectString 填充得到数据库连接字符串
func ReadMysqlString() {
	fmt.Sscanf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		Info.DB.Username, Info.DB.Password, Info.DB.Host, Info.DB.Port, Info.DB.Database,
		Info.DB.Charset, Info.DB.ParseTime, Info.DB.Loc)
	//log.Println(arg)
}

var Reader *viper.Viper

func Init() {
	Reader = viper.New()
	path, _ := os.Getwd()
	Reader.AddConfigPath(path + "./config")
	Reader.SetConfigName("config")
	Reader.SetConfigType("yaml")
	err := Reader.ReadInConfig() // 查找并读取配置文件
	if err != nil {              // 处理读取配置文件的错误
		logrus.Error("Read config file failed: %s \n", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Info("no error in config file")
		} else {
			logrus.Error("found error in config file\n", ok)
		}
	}
}
