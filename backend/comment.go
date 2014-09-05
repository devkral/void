package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo/bson"
)

type Commment struct {
	Id   bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Text string
	Date string
	User *bson.ObjectId
	Type string
}

type CommentResource struct{}

func (r CommentResource) Register(wsContainer *restful.Container) {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path("/rest/comments")

	ws.Route(ws.GET("/").Filter(authFilter).To(r.getComments))
	ws.Route(ws.POST("/").Filter(authFilter).To(r.createComment))
	ws.Route(ws.PUT("/{entry}").Filter(authFilter).To(r.editComment))
	ws.Route(ws.DELETE("/{entry}").Filter(authFilter).To(r.deleteComment))
    wsContainer.Add(ws)
}

func (r CommentResource) getComments(req *restful.Request, resp *restful.Response) {
	//TODO:implement
}

func (r CommentResource) createComment(req *restful.Request, resp *restful.Response) {
	//TODO:implement
}

func (r CommentResource) editComment(req *restful.Request, resp *restful.Response) {
	//TODO:implement
}

func (r CommentResource) deleteComment(req *restful.Request, resp *restful.Response) {
	//TODO:implement
}
