# 根据当前操作系统，将openvscode中的内容解压third-party中
OS = $(shell uname)
ARCH = $(shell uname -m)

ifeq (${OS}, Darwin)
	OS = darwin
else ifeq (${OS}, Linux)
	OS = linux
else
	$(error Unsupported OS ${OS})
endif

ifeq (${ARCH}, x86_64)
	ARCH = x64
endif

VSCODE_SERVER = openvscode-server-v1.67.0-${OS}-${ARCH}

third_party:
	@if [ -d third_party/${VSCODE_SERVER} ]; then echo "vscode server is ready"; else mkdir -p third_party/${VSCODE_SERVER} && curl https://s-public-packages.oss-cn-hangzhou.aliyuncs.com/openvscode-server/${VSCODE_SERVER}.tar.gz -o third_party/${VSCODE_SERVER}.tar.gz && tar zxvf third_party/${VSCODE_SERVER}.tar.gz -C third_party; fi

release:
	go clean
	rm -rf ./target/*
	GOOS=linux GOARCH=amd64 go build -o target/vscode-server src/main.go
	cp config/fc.yaml target/config.yaml

clean:
	go clean
	rm -rf ./target/*