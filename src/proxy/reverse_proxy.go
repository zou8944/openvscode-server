package proxy

import (
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
		url, err := url2.Parse("http://" + sp.VscodeServer.Host + ":" + sp.VscodeServer.Port)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		sp.Proxy = httputil.NewSingleHostReverseProxy(url)
		sp.VscodeServer.Init()
	}
}

func (sp *ServerProxy) shutdownHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sp.VscodeServer.Shutdown()
		if sp.VscodeServer.Err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (sp *ServerProxy) processHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sp.Proxy.ServeHTTP(w, r)
	}
}

func StartProxy(vscodeServer *vscode.Server) error {
	sp := &ServerProxy{
		VscodeServer: vscodeServer,
	}
	http.HandleFunc("/initialize", sp.initializeHandler())
	http.HandleFunc("/pre-stop", sp.shutdownHandler())
	http.HandleFunc("/", sp.processHandler())

	s := &http.Server{
		Addr:           "127.0.0.1:8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
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
