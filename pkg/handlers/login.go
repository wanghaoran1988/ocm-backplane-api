package handlers

import "net/http"


// LoginHandler will return back the cluster url for client
func LoginHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

	}
	return http.HandlerFunc(fn)
}
