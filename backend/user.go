package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo/bson"
)

type User struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Email        string
	Name         string
	Organization string
}

func LoadUserByName(name string) (*User, error) {
	//TODO: implement
	return nil, nil
}

func (u *User) Authenticate(pw string) bool {
	//TODO: implement
	return false
}

type UserResource struct{}

func (r UserResource) Register(wsContainer *restful.Container) {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path("/rest/users")

	ws.Route(ws.GET("/").Filter(authFilter).To(r.getUsers))
	ws.Route(ws.POST("/").Filter(authFilter).To(r.createUser))
	ws.Route(ws.PUT("/{entry}").Filter(authFilter).To(r.editUser))
	ws.Route(ws.DELETE("/{entry}").Filter(authFilter).To(r.deleteUser))
}

func (r UserResource) getUsers(req *restful.Request, resp *restful.Response) {
	//TODO:implement
}

func (r UserResource) createUser(req *restful.Request, resp *restful.Response) {
	//TODO:implement
}

func (r UserResource) editUser(req *restful.Request, resp *restful.Response) {
	//TODO:implement
}

func (r UserResource) deleteUser(req *restful.Request, resp *restful.Response) {
	//TODO:implement
}
