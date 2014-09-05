package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo/bson"
    "net/http"
)

type Comment struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Text     string
	Date     string
	User     bson.ObjectId `json:"user"`
	Type     string
	Building bson.ObjectId `json:"building"`
}

type CommentWrapper struct {
    Comment *Comment
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
    reqUser := getRequestUser(req)
    if reqUser == nil {
        resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
        return
    }
	//TODO:implement
}

func (r CommentResource) createComment(req *restful.Request, resp *restful.Response) {
    reqUser := getRequestUser(req)
    if reqUser == nil {
        resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
        return
    }
	//TODO:implement
    co := new(CommentWrapper)
    err := req.ReadEntity(co)
    if err == nil {
        co.Comment.Save()
        resp.WriteEntity(co)
    } else {
        resp.WriteErrorString(http.StatusBadRequest, "Your comment makes no sense at all!")
    }
}

func (r CommentResource) editComment(req *restful.Request, resp *restful.Response) {
    reqUser := getRequestUser(req)
    if reqUser == nil {
        resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
        return
    }
	//TODO:implement
}

func (r CommentResource) deleteComment(req *restful.Request, resp *restful.Response) {
    reqUser := getRequestUser(req)
    if reqUser == nil {
        resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
        return
    }
	//TODO:implement
}

func (c *Comment) Save() error {
    if !c.Id.Valid() {
        c.Id = bson.NewObjectId()
    }
    _, err := mongo.DB("void").C("comments").UpsertId(c.Id, c)
    return err
}
