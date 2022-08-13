package vscode

import (
	"aliyun/serverless/webide-server/src/context"
	"github.com/spf13/viper"
	"os"
	"path"
	"runtime"
	"testing"
)

func TestServer_Init(t *testing.T) {
	server := createServer(t)
	if server.Err != nil {
		t.Error(server.Err)
	}
	server.Init()
	if server.Err != nil {
		t.Error(server.Err)
	}
	server.Shutdown()
	if server.Err != nil {
		t.Error(server.Err)
	}
}

func chdir() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func createServer(t *testing.T) *Server {
	chdir()
	ctx, err := context.CreateFromEnv()
	if err != nil {
		t.Error(err)
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile("config/test.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		t.Error(err)
	}
	// 创建Server
	server, err := NewServer(ctx)

	if err != nil {
		t.Error(err)
	}
	return server
}
