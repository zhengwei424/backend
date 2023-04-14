package tools

import (
	"bufio"
	"encoding/json"
	"os"
)

// App 后端启动配置
type App struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// Redis redis配置，存储验证码等
type Redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Mysql struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Config 全局综合配置
type Config struct {
	AppConfig   *App   `json:"app"`
	RedisConfig *Redis `json:"redis"`
	MysqlConfig *Mysql `json:"mysql"`
}

// GetConfig 解析全局配置，并返回
func GetConfig() *Config {
	var _cfg = new(Config)
	// 固定解析app.json文件，全局配置都在这里
	file, err := os.Open("globalConfig/dbConfig.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(_cfg); err != nil {
		panic(err)
	}
	return _cfg
}
