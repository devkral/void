package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo/bson"
    "net/http"
)

type Invitation struct {
	Id    bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Email string
}

type InvitationResource struct{}

func (r InvitationResource) Register(wsContainer *restful.Container) {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path("/rest/invitations")

	ws.Route(ws.GET("/{entry}").Filter(authFilter).To(r.getInvitation))
	ws.Route(ws.POST("/").Filter(authFilter).To(r.createInvitation))
	wsContainer.Add(ws)
}

func (r InvitationResource) getInvitation(req *restful.Request, resp *restful.Response) {
	//TODO:implement
}

func (r InvitationResource) createInvitation(req *restful.Request, resp *restful.Response) {
    reqUser := getRequestUser(req)
    if reqUser == nil {
        resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
        return
    }
	//TODO:implement
}
