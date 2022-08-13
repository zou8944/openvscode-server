package vscode

import (
	"aliyun/serverless/webide-server/src/compress"
	"aliyun/serverless/webide-server/src/context"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

/**
这个文件里我需要做的事情
1.启动服务器
	- 从oss加载文件到本地目录
	- 通过命令行启动vscode服务器
2.关闭服务器
	- 下线服务器
	- 将本地目录打包，上传到oss
需要的资料
- vscode服务器的使用方式
- oss的使用方式
- 命令行加载的使用方式
*/

type Server struct {
	Host              string
	Port              string
	vscodeBinaryPath  string // vscode的可执行文件存放位置
	DataDir           string // 存储用户数据、服务数据、扩展信息数据的地方
	WorkspaceDir      string // 工作空间数据
	OssBucket         string
	OssPath4Data      string
	OssPath4Workspace string
	OssClient         *oss.Client
	cmd               *exec.Cmd
	Err               error
}

func NewServer(ctx *context.Context) (*Server, error) {
	ossClient, err := oss.New("oss-"+ctx.Region+".aliyuncs.com", ctx.AccessKeyId, ctx.AccessKeySecret)
	if err != nil {
		return nil, err
	}

	s := &Server{
		Host:              viper.GetString("vscode.host"),
		Port:              viper.GetString("vscode.port"),
		vscodeBinaryPath:  viper.GetString("vscode.binaryPath"),
		DataDir:           viper.GetString("vscode.dataDir"),
		WorkspaceDir:      viper.GetString("vscode.workspaceDir"),
		OssBucket:         viper.GetString("oss.bucket"),
		OssPath4Data:      viper.GetString("oss.dataPath"),
		OssPath4Workspace: viper.GetString("oss.workspacePath"),
		OssClient:         ossClient,
	}
	// 以环境变量为准
	ossBucket := os.Getenv("OSS_BUCKET_NAME")
	if ossBucket != "" {
		s.OssBucket = ossBucket
	}

	log.Printf("VsCode服务器创建成功: %+v", s)

	return s, nil
}

func (s *Server) loadData() {
	if s.Err != nil {
		return
	}
	bucket, err := s.OssClient.Bucket(s.OssBucket)
	if err != nil {
		s.Err = err
		return
	}

	// 用户数据
	dataObject, err := bucket.GetObject(s.OssPath4Data)
	if err != nil {
		// 首次是读取不到对象的
		if err.(oss.ServiceError).StatusCode == http.StatusNotFound {
			log.Println("OSS找不到对应object", s.OssPath4Data)
		} else {
			s.Err = err
			return
		}
	} else {
		dataObjectBytes, err := ioutil.ReadAll(dataObject)
		if err != nil {
			s.Err = err
			return
		}
		err = compress.DeCompress(s.DataDir, dataObjectBytes)
		if err != nil {
			s.Err = err
			return
		}
	}
	// 空间数据
	workspaceObject, err := bucket.GetObject(s.OssPath4Workspace)
	if err != nil {
		if err.(oss.ServiceError).StatusCode == http.StatusNotFound {
			log.Println("OSS找不到对应object", s.OssPath4Data)
		} else {
			s.Err = err
			return
		}
	} else {
		dataObjectBytes, err := ioutil.ReadAll(workspaceObject)
		if err != nil {
			s.Err = err
			return
		}
		err = compress.DeCompress(s.WorkspaceDir, dataObjectBytes)
		if err != nil {
			s.Err = err
			return
		}
	}
}

func (s *Server) saveData() {
	if s.Err != nil {
		return
	}
	bucket, err := s.OssClient.Bucket(s.OssBucket)
	if err != nil {
		s.Err = err
		return
	}
	// 目录压缩成文件
	file, err := os.CreateTemp("", s.OssPath4Data)
	err = compress.Compress(file, s.DataDir)
	err = bucket.PutObjectFromFile(s.OssPath4Data, file.Name())
	err = file.Close()

	file, err = os.CreateTemp("", s.OssPath4Workspace)
	err = compress.Compress(file, s.WorkspaceDir)
	err = bucket.PutObjectFromFile(s.OssPath4Workspace, file.Name())
	err = file.Close()
	if err != nil {
		s.Err = err
		return
	}
}

func (s *Server) Init() {
	if s.Err != nil {
		return
	}
	s.loadData()
	if s.Err != nil {
		return
	}
	userDataDir := filepath.Join(s.DataDir, "user-data")
	serverDataDir := filepath.Join(s.DataDir, "server-data")
	extensionsDir := filepath.Join(s.DataDir, "extensions")
	s.cmd = exec.Command(
		s.vscodeBinaryPath,
		"--host="+s.Host,
		"--port="+s.Port,
		"--user-data-dir="+userDataDir,
		"--server-data-dir="+serverDataDir,
		"--extensions-dir="+extensionsDir,
		"--without-connection-token",
		"--start-server",
		"--telemetry-level=off",
	)
	log.Println(s.cmd.Args)
	err := s.cmd.Start()
	if err != nil {
		s.Err = err
		return
	}
	// 等待服务启动
	for {
		if _, err := net.Dial("tcp", s.Host+":"+s.Port); err == nil {
			log.Println("vscode server 启动成功")
			break
		} else {
			log.Println("等待 vscode server 启动")
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func (s *Server) Shutdown() {
	if s.Err != nil {
		return
	}
	s.saveData()
	err := s.cmd.Process.Release()
	if err != nil {
		s.Err = err
	}
}
