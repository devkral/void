package main

import (
	"encoding/base64"
	"github.com/emicklei/go-restful"
	"strings"
)

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
	u, err = LoadUserByName(authparams[0])
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
