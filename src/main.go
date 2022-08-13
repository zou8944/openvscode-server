package main

import (
	"aliyun/serverless/webide-server/src/proxy"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

func main() {
	pwd, err := os.Executable()
	if err != nil {
		log.Fatalln("获取当前运行目录失败", err)
	}
	// 在fc上，配置文件会在config.yaml总；直接运行时，在config/dev.yaml中
	stage := os.Getenv("STAGE")
	var configDir string
	if stage == "dev" {
		pwd, err = os.Getwd()
		configDir = filepath.Join(filepath.Dir(pwd), "serverless-vscode/config/dev.yaml")
	} else {
		configDir = filepath.Join(filepath.Dir(pwd), "config.yaml")
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(configDir)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("加载配置文件出错, ", err)
	}
	err = proxy.StartProxy()
	if err != nil {
		log.Fatalln("反向代理启动失败", err)
	}
}
