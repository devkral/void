/*
 *   © 2014 Daniel 'grindhold' Brendle and Contributors
 *
 *   This file is part of Void.
 *
 *   Void is free software: you can redistribute it and/or
 *   modify it under the terms of the GNU Affero General Public License
 *   as published by the Free Software Foundation, either
 *   version 3 of the License, or (at your option) any later
 *   version.
 *
 *   Void is distributed in the hope that it will be
 *   useful, but WITHOUT ANY WARRANTY; without even the implied
 +   warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
 *   PURPOSE. See the GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public
 *   License along with Void.
 *   If not, see http://www.gnu.org/licenses/.
*/

package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"time"
)

type Comment struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Text     string
	Date     string
	User     bson.ObjectId `json:"user"`
  Logcomment bool
	Building bson.ObjectId `json:"building"`
}

type CommentWrapper struct {
	Comment *Comment
}

type CommentsWrapper struct {
	Comments []*Comment
	Users    []*User
}

func (c *Comment) Save() error {
	if !c.Id.Valid() {
		c.Id = bson.NewObjectId()
		c.Date = time.Now().Format(time.RFC3339)
		if building, err := LoadBuildingById(c.Building); err != nil {
			log.Println("could not save comment to nonexistent building " + c.Building.Hex())
		} else {
			building.AddComment(c)
		}
	}
	_, err := mongo.DB("void").C("comments").UpsertId(c.Id, c)
	return err
}

func (c *Comment) Delete() error {
	b, err := LoadBuildingById(c.Building)
	if err == nil {
		b.RemoveComment(c)
	}
	return mongo.DB("void").C("comments").RemoveId(c.Id)
}

type CommentResource struct{}

func (r CommentResource) Register(wsContainer *restful.Container) {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path("/rest/comments")

	ws.Route(ws.GET("/").Filter(authFilter).To(r.getComments))
	ws.Route(ws.POST("/").Filter(authFilter).To(r.createComment))
	ws.Route(ws.DELETE("/{entry}").Filter(authFilter).To(r.deleteComment))
	wsContainer.Add(ws)
}

func (r CommentResource) getComments(req *restful.Request, resp *restful.Response) {
	reqUser := getRequestUser(req)
	if reqUser == nil {
		resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
		return
	}
	if arr, ok := req.Request.URL.Query()["ids[]"]; ok {
		ret := make([]*Comment, 0)
		users := make([]*User, 0)
		handled_users := make(map[bson.ObjectId]bool)
		for i := range arr {
			if c, err := LoadCommentById(bson.ObjectIdHex(arr[i])); err == nil {
				ret = append(ret, c)
				if _, ok := handled_users[c.User]; !ok {
					u, err := LoadUserById(c.User)
					if err != nil {
						continue
					}
					handled_users[c.User] = true
					users = append(users, u)
				}
			}
		}
		cw := new(CommentsWrapper)
		cw.Comments = ret
		cw.Users = users
		resp.WriteEntity(cw)
	} else {
		resp.WriteErrorString(http.StatusBadRequest, "need ids[]")
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
		co.Comment.Logcomment = false
		err2 := co.Comment.Save()
		if err2 != nil {
			resp.WriteErrorString(http.StatusInternalServerError, err2.Error())
			return
		}
		resp.WriteEntity(co)
	} else {
		resp.WriteErrorString(http.StatusBadRequest, "Your comment makes no sense at all!")
	}
}

/*func (r CommentResource) editComment(req *restful.Request, resp *restful.Response) {
    reqUser := getRequestUser(req)
    if reqUser == nil {
        resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
        return
    }
	//TODO:implement
}*/

func (r CommentResource) deleteComment(req *restful.Request, resp *restful.Response) {
	reqUser := getRequestUser(req)
	if reqUser == nil {
		resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
		return
	}
	if c, err := LoadCommentById(bson.ObjectIdHex(req.PathParameter("entry"))); err == nil {
		if c.User == reqUser.Id {
			c.Delete()
			resp.WriteEntity(Empty{})
		}
	}
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
