package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"flag"
	"github.com/elazarl/goproxy"

	log "github.com/liudanking/goutil/logutil"

	kcp "github.com/xtaci/kcp-go/v5"
)

func main() {
	var (
		listenAddr   string
		cert         string
		key          string
		auth         string
		verbose      bool
		printVersion bool
	)
	flag.StringVar(&listenAddr, "l", ":443", "listen addr (udp port only)")
	flag.StringVar(&cert, "cert", "", "cert path")
	flag.StringVar(&key, "key", "", "key path")
	flag.StringVar(&auth, "auth", "quic-proxy:Go!", "basic auth, format: username:password")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.BoolVar(&printVersion, "V", false, "print version")
	flag.Parse()

	if printVersion {
		fmt.Fprintf(os.Stdout, "Quic Server %s (%s %s/%s)\n",
			"2.0 2022-02-09", runtime.Version(), runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	log.Info("%v", verbose)
	if cert == "" || key == "" {
		log.Error("cert and key can't by empty")
		return
	}

	parts := strings.Split(auth, ":")
	if len(parts) != 2 {
		log.Error("auth param invalid")
		return
	}
	username, password := parts[0], parts[1]

	listener, err := kcp.ListenWithOptions(listenAddr, nil, 10, 3)

	if err != nil {
		log.Error("listen failed:%v", err)
		return
	}

	proxy := goproxy.NewProxyHttpServer()
	ProxyBasicAuth(proxy, func(u, p string) bool {
		return u == username && p == password
	})
	proxy.Verbose = verbose
	server := &http.Server{Addr: listenAddr, Handler: proxy}
	log.Info("start serving %v", listenAddr)
	log.Error("serve error:%v", server.Serve(listener))

}
