# Serverless VSCode

来源于示例代码，当前功能和示例代码没有任何区别，主要目的是跑通。

## 前置知识

- Go语言不必多说
- 阿里云函数计算
  - http函数的构建方式
  - 函数生命周期回调
- Serverless Devs
  - 使用方式
  - fc组件
  - layer-fc插件
- openvscode-server命令行的使用方式

项目实现本身没有什么难度，需要花费很多时间的是函数计算的部署和调试。

## 基本原理

- 项目本身构建了一个简单的反向代理服务器，总共暴露三个处理器
  - /initialize 用于加载用户数据然后启动openvscode-server服务。设置在函数实例启动时执行
  - /pre-stop 用于保存用户数据并上传oss。设置在函数实例关闭前执行
  - /* 其它所有请求，直接转发给openvscode-server服务

## 改动

与原项目有几点改动

- 移除不必要的内容，原项目可以直接发布到Serverless Devs商店，我们这里不用，把那些多余的删除
- 原项目的fc环境下配置文件加载其实是不生效的，能工作是因为使用了默认值
- 原项目代码比较混乱，我选择自己重写，当然我的代码也好不到哪里去，但考虑到我是个go新手就原谅我吧
- 原项目配置是需要输入的，为了方便审核人员测试，我将其中的账号固定为我自己私人的RAM账号（十天后将被轮转）

## 本地运行

设置环境变量

```shell
ALIYUN_ACCESS_KEY_ID=xxx
ALIYUN_ACCESS_KEY_SECRET=xxx
ALIYUN_OSS_REGION=cn-xxx
STAGE=dev
```

然后 `go run src/main.go`

## 部署

cd到项目根目录，然后

`s deploy`

或者在函数计算控制台直接从code space创建，不需要填写任何参数