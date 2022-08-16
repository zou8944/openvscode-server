package proxy

import (
	"aliyun/serverless/webide-server/src/context"
	"aliyun/serverless/webide-server/src/vscode"
	"log"
	"net/http"
	"net/http/httputil"
	url2 "net/url"
	"time"
)

type ServerProxy struct {
	VscodeServer *vscode.Server
	Proxy        *httputil.ReverseProxy
}

func (sp *ServerProxy) initializeHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("进入初始化")
		// 直接从环境变量中取得
		ctx, err := context.CreateFromEnv()
		if err != nil {
			log.Fatalln("Context加载失败", err)
		}
		vscodeServer, err := vscode.NewServer(ctx)
		if err != nil {
			log.Fatalln("Vs code server创建失败", err)
		}
		url, err := url2.Parse("http://" + vscodeServer.Host + ":" + vscodeServer.Port)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sp.Proxy = httputil.NewSingleHostReverseProxy(url)
		sp.VscodeServer = vscodeServer
		sp.VscodeServer.Init()
		log.Println("初始化完成")
	}
}

func (sp *ServerProxy) shutdownHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("预关闭")
		sp.VscodeServer.Shutdown()
		if sp.VscodeServer.Err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("预关闭完成")
	}
}

func (sp *ServerProxy) processHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("一般请求进来了: %+v\n", r)
		sp.Proxy.ServeHTTP(w, r)
	}
}

func StartProxy() error {
	sp := &ServerProxy{}
	http.HandleFunc("/initialize", sp.initializeHandler())
	http.HandleFunc("/pre-stop", sp.shutdownHandler())
	http.HandleFunc("/", sp.processHandler())

	s := &http.Server{
		Addr:           ":9000",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    5 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("服务启动成功: " + s.Addr)
	err := s.ListenAndServe()
	if err != nil {
		return err
	}

	log.Println("服务启动成功: " + s.Addr)
	return nil
}
