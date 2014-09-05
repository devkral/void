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

type CommentsWrapper struct {
    Comments []*Comment
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
    c, err := LoadComments()
    if err == nil {
        cw := new(CommentsWrapper)
        cw.Comments = c
        resp.WriteEntity(cw)
    } else {
        resp.WriteErrorString(http.StatusInternalServerError, "Nothing found")
    }
}

func (r CommentResource) createComment(req *restful.Request, resp *restful.Response) {
    reqUser := getRequestUser(req)
    if reqUser == nil {
        resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
        return
    }
    co := new(CommentWrapper)
    err := req.ReadEntity(co)
    if err == nil {
        err2 := co.Comment.Save()
        if err2 != nil {
            resp.WriteErrorString(http.StatusInternalServerError, err2.Error())
        }
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

func LoadCommentById(id bson.ObjectId) (*Comment, error) {
    x := new(Comment)
    err := mongo.DB("void").C("comments").Find(bson.M{"_id": id}).One(x)
    return x, err
}

func LoadComments() ([]*Comment, error) {
    x := make([]*Comment, 0)
    err := mongo.DB("void").C("buildings").Find(bson.M{}).All(&x)
    return x, err
}

