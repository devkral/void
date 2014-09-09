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
	"encoding/base64"
	"github.com/emicklei/go-restful"
	"strings"
)

type AuthenticationResult struct {
	Valid bool
}

func authFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	var authstr_bytes []byte
	var authstr string
	var authparams []string
	var u *User
	var err error

	authstr = req.Request.Header.Get("Authorization")
	if len(authstr) > 6 {
		authstr = authstr[6:]
	} else {
		goto noauth
	}
	authstr_bytes, err = base64.StdEncoding.DecodeString(authstr)
	if err != nil {
		goto noauth
	}
	authparams = strings.SplitN(string(authstr_bytes), ":", 2)
	if len(authparams) != 2 {
		goto noauth
	}
	u, err = LoadUserByEmail(authparams[0])
	if err != nil {
		goto noauth
	}
	if u.Authenticate(authparams[1]) {
		setRequestUser(req, u)
		defer unsetRequestUser(req)
	}

noauth:
	chain.ProcessFilter(req, resp)
}

type AuthResource struct{}

func (r AuthResource) Register(wsContainer *restful.Container) {
	requestUserMap = map[*restful.Request]*User{}

	ws := new(restful.WebService)
	ws.Produces(restful.MIME_JSON)
	ws.Path("/auth")
	ws.Route(ws.GET("/").Filter(authFilter).To(r.authHandler))
	wsContainer.Add(ws)
}

func (r AuthResource) authHandler(req *restful.Request, resp *restful.Response) {
	success := getRequestUser(req) != nil
	resp.WriteEntity(AuthenticationResult{Valid: success})
	return
}

var requestUserMap map[*restful.Request]*User

func setRequestUser(req *restful.Request, u *User) {
	requestUserMap[req] = u
}

func unsetRequestUser(req *restful.Request) {
	delete(requestUserMap, req)
}

func getRequestUser(req *restful.Request) *User {
	if u, ok := requestUserMap[req]; ok {
		return u
	} else {
		return nil
	}
}
