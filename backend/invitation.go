package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo/bson"
)

type Invitation struct {
	Id bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
}

type InvitationResource struct{}

func (r InvitationResource) Register(wsContainer *restful.Container) {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path("/rest/invitations")

	ws.Route(ws.GET("/").Filter(authFilter).To(r.getInvitations))
	ws.Route(ws.POST("/").Filter(authFilter).To(r.createInvitation))
	ws.Route(ws.PUT("/{entry}").Filter(authFilter).To(r.editInvitation))
}

func (r InvitationResource) getInvitations(req *restful.Request, resp *restful.Response) {
	//TODO:implement
}

func (r InvitationResource) createInvitation(req *restful.Request, resp *restful.Response) {
	//TODO:implement
}

func (r InvitationResource) editInvitation(req *restful.Request, resp *restful.Response) {
	//TODO:implement
}
