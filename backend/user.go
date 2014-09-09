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
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo/bson"
	"math/rand"
	"net/http"
)

type UserWrapper struct {
	User *User
}

type UsersWrapper struct {
	Users []*User
}

type User struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Email        string
	Organization string

	IPassword string `json:"-"`
	ISalt     int64  `json:"-"`
	Password  string
}

func LoadUserByEmail(email string) (*User, error) {
	u := new(User)
	err := mongo.DB("void").C("users").Find(bson.M{"email": email}).One(&u)
	return u, err
}

func LoadUserById(id bson.ObjectId) (*User, error) {
	u := new(User)
	err := mongo.DB("void").C("users").Find(bson.M{"_id": id}).One(&u)
	return u, err
}

func InitializeAdmin() {
	if _, err := LoadUserByEmail("admin@nonexistent.invalid"); err != nil {
		admin := new(User)
		admin.Email = "admin@nonexistent.invalid"
		admin.Organization = "adminorga"
		admin.SetPassword("admin")
		admin.Save()
	}
}

func (u *User) Authenticate(pw string) bool {
	s := make([]byte, 8)
	binary.LittleEndian.PutUint64(s, uint64(u.ISalt))
	pw_bytes := []byte(pw)
	sum := sha512.Sum512([]byte(append(pw_bytes, s...)))
	return u.IPassword == hex.EncodeToString(sum[0:64])
}

func (u *User) SetPassword(pw string) {
	u.ISalt = rand.Int63()
	s := make([]byte, 8)
	binary.LittleEndian.PutUint64(s, uint64(u.ISalt))
	pw_bytes := []byte(pw)
	sum := sha512.Sum512([]byte(append(pw_bytes, s...)))
	u.IPassword = hex.EncodeToString(sum[0:64])
	u.Password = "" //Ensure the password will not be returned in PUT-JSON
}

func (u *User) Save() error {
	if !u.Id.Valid() {
		u.Id = bson.NewObjectId()
	}
	_, err := mongo.DB("void").C("users").UpsertId(u.Id, u)
	return err
}

type UserResource struct{}

func (r UserResource) Register(wsContainer *restful.Container) {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path("/rest/users")

	ws.Route(ws.GET("/").Filter(authFilter).To(r.getUsers))
	ws.Route(ws.PUT("/{entry}").Filter(authFilter).To(r.editUser))
	ws.Route(ws.DELETE("/{entry}").Filter(authFilter).To(r.deleteUser))
	wsContainer.Add(ws)
}

func (r UserResource) getUsers(req *restful.Request, resp *restful.Response) {
	reqUser := getRequestUser(req)
	if reqUser == nil {
		resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
		return
	}
	if arr, ok := req.Request.URL.Query()["email"]; ok && len(arr) == 1 {
		user, err := LoadUserByEmail(arr[0])
		if err != nil {
			resp.WriteErrorString(http.StatusNotFound, "no such user")
		} else {
			uw := new(UsersWrapper)
			uw.Users = []*User{user}
			resp.WriteEntity(uw)
		}
	}
}

func (r UserResource) editUser(req *restful.Request, resp *restful.Response) {
	reqUser := getRequestUser(req)
	if reqUser == nil {
		resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
		return
	}
	//TODO:implement
}

func (r UserResource) deleteUser(req *restful.Request, resp *restful.Response) {
	reqUser := getRequestUser(req)
	if reqUser == nil {
		resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
		return
	}
	//TODO:implement
}
