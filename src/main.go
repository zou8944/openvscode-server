package main

import (
	"aliyun/serverless/webide-server/src/context"
	"aliyun/serverless/webide-server/src/proxy"
	"aliyun/serverless/webide-server/src/vscode"
	"github.com/spf13/viper"
	"log"
)

func main() {
	/**
	还缺少的内容
	1. 根据不同的环境，加载不同的配置文件到viper
	2. 根据不同的环境，context需要从不同的位置提取参数
	3. 缺少阿里云函数计算应用的配置，需要一键部署
	4. 缺少明确的项目发布方式：开发环境、测试环境、线上环境分别应该怎么发布？
	5. 缺少测试
	*/
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config/fc.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("加载配置文件出错, ", err)
	}
	ctx, err := context.CreateFromEnv()
	if err != nil {
		log.Fatalln("Context加载失败", err)
	}
	vscodeServer, err := vscode.NewServer(ctx)
	err = proxy.StartProxy(vscodeServer)
	if err != nil {
		log.Fatalln("反向代理启动失败", err)
	}
}
