package login

import (
	"github.com/emicklei/go-restful"
)

// WebService creates a new service that can handle REST requests for User resources.
func WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/backplane").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/login").To(login))

	return ws
}

// TODO get the cluster id/name, user/group info from the request headers, and create related service account, and rolebindings ?
// questions:
// 1. how to verify a ocm user can run ocm backplen login <cluster_id>, if there is no scope available for this user, then failed?
// from the cluster info, we can know it's a aws/gce
func login(request *restful.Request, response *restful.Response) {
	response.Write([]byte("login succeed"))
}
