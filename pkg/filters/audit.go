package filters

import (
	"net/http"
)

//TODO
func WithAudit(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(w, req)
		//_ := w.Header().Get("Audit-Id")
	})
}
