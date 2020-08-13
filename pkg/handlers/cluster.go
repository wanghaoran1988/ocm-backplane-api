package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/wanghaoran1988/ocm-backplane-api/pkg/filters"
	"github.com/wanghaoran1988/ocm-backplane-api/pkg/ocm"
	"github.com/wanghaoran1988/ocm-backplane-api/pkg/proxy"
)

var (
	defaultAPIPrefix = "/backplane/cluster/"
)

type ClusterApiHandler struct {
	ocm.KubecfgGetter
}

func NewClusterApiHandler(configGetter ocm.KubecfgGetter) http.Handler {
	return &ClusterApiHandler{
		configGetter,
	}
}

func (h *ClusterApiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	kubeconfig := h.KubecfgGetter.GetKubeConfig(id)
	if kubeconfig == nil {
		http.NotFound(w,req)
	}
	proxyApi := defaultAPIPrefix + id + "/"

	proxy, err := proxy.NewProxy(proxyApi, kubeconfig, time.Duration(0))
	if err != nil {
		log.Printf("failed to create proxy handler: %v", err)
	}
	proxy = filters.WithAudit(proxy)

	proxy.ServeHTTP(w, req)
}
