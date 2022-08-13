package context

import (
	"fmt"
	"net/http"
	"os"
)

type Context struct {
	AccessKeyId     string
	AccessKeySecret string
	Region          string
}

// CreateFromHttpRequest 从请求或环境变量中创建上下文对象
func CreateFromHttpRequest(request *http.Request) (*Context, error) {
	ctx := &Context{
		AccessKeyId:     request.Header.Get("x-fc-access-key-id"),
		AccessKeySecret: request.Header.Get("x-fc-access-key-secret"),
		Region:          request.Header.Get("x-fc-region"),
	}
	if ctx.AccessKeyId == "" {
		return nil, fmt.Errorf("can not get access key id from fc runtime. Please make sure you already granted OSS permission to your FC function")
	}
	if ctx.AccessKeySecret == "" {
		return nil, fmt.Errorf("can not get access key secret from fc runtime. Please make sure you already granted OSS permission to your FC function")
	}
	if ctx.Region == "" {
		return nil, fmt.Errorf("can not get region from fc runtime. Please make sure you already granted OSS permission to your FC function")
	}
	return ctx, nil
}

func CreateFromEnv() (*Context, error) {
	ctx := &Context{
		AccessKeyId:     os.Getenv("ALIYUN_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("ALIYUN_ACCESS_KEY_SECRET"),
		Region:          os.Getenv("ALIYUN_OSS_REGION"),
	}
	if ctx.AccessKeyId == "" {
		return nil, fmt.Errorf("can not get access key id from environment variable")
	}
	if ctx.AccessKeySecret == "" {
		return nil, fmt.Errorf("can not get access key secret from environment variable")
	}
	if ctx.Region == "" {
		return nil, fmt.Errorf("can not get region from environment variable")
	}

	return ctx, nil
}
