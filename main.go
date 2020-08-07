package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/wanghaoran1988/ocm-backplane-api/pkg/proxy"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultPort         = 8001
	defaultStaticPrefix = "/static/"
	defaultAPIPrefix    = "/backplane/cluster/"
	defaultAddress      = "127.0.0.1"
)

func main() {

	configDir := flag.String("configdir", "", "the directory that contains all the kubeconfig")

	flag.Parse()

	var configFiles []string
	err := filepath.Walk(*configDir, func(path string, info os.FileInfo, err error) error {
		fmt.Printf("visited file or dir: %q\n", path)
		if !info.IsDir() {
			configFiles = append(configFiles, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	var kubeconfigs []*rest.Config
	mux := http.NewServeMux()
	for _, file := range configFiles {
		fmt.Printf("load kubeconfig %s\n", file)
		// use the current context in kubeconfig
		clientConfig, err := clientcmd.BuildConfigFromFlags("", file)
		if err != nil {
			panic(err.Error())
		}
		kubeconfigs = append(kubeconfigs, clientConfig)
		proxyApi := defaultAPIPrefix + filepath.Base(file) +"/"
		fmt.Printf("Adding proxy api: %s\n", proxyApi)
		proxy, err := proxy.NewProxy(proxyApi, clientConfig, time.Duration(0))
		if err != nil {
			log.Printf("failed to create proxy handler: %v", err)
		}
		mux.Handle(proxyApi, proxy)

	}

	server, err := proxy.NewServer(mux)

	if err != nil {

		log.Fatalf("failed to create proxy server: %v", err)
	}

	// Separate listening from serving so we can report the bound port
	// when it is chosen by os (eg: port == 0)
	var l net.Listener

	l, err = server.Listen(defaultAddress, defaultPort)

	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(server.ServeOnListener(l))
}
