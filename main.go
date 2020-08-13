package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/wanghaoran1988/ocm-backplane-api/pkg/handlers"
	"github.com/wanghaoran1988/ocm-backplane-api/pkg/ocm"
	"github.com/wanghaoran1988/ocm-backplane-api/pkg/proxy"
	"log"
	"net"
)

const (
	defaultPort    = 8001
	defaultAddress = "127.0.0.1"
)

func main() {

	configDir := flag.String("configdir", "", "the directory that contains all the kubeconfig")

	flag.Parse()

	configGetter := ocm.NewConfigFileGetter(*configDir)
	err := configGetter.Init()
	if err != nil {
		log.Fatalf("failed to init kube config getter: %v", err)
	}
	r := mux.NewRouter()
	clusterHandler := handlers.NewClusterApiHandler(configGetter)
	r.PathPrefix("/backplane/cluster/{id}/").Handler(clusterHandler)
	r.Handle("/backplane/login", handlers.LoginHandler())
	server, err := proxy.NewServer(r)

	if err != nil {

		log.Fatalf("failed to create proxy server: %v", err)
	}

	var l net.Listener

	l, err = server.Listen(defaultAddress, defaultPort)

	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(server.ServeOnListener(l))
}
