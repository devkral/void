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
	"github.com/grindhold/gominatim"
	"io/ioutil"
	"labix.org/v2/mgo"
	"log"
	"net/http"
)

var mongo *mgo.Session

func main() {
	log.Println("Entering the void.")
	log.Println("\tEstablishing connection to mongo DB...")
	mng, err := mgo.Dial("localhost")
	mongo = mng
	if err != nil {
		log.Fatal("Could not connect to mongodb!")
	}
	// Prepare REST-backend
	log.Println("\tInitializing REST-Backend...")
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	BuildingResource{}.Register(wsContainer)
	UserResource{}.Register(wsContainer)
	CommentResource{}.Register(wsContainer)
	InvitationResource{}.Register(wsContainer)
	AuthResource{}.Register(wsContainer)
	ViewResource{}.Register(wsContainer)
	StaticResource{}.Register(wsContainer)

	InitializeAdmin()

	//Initialize Gominatim
	log.Println("\tInitializing gominatim")
	gominatim.SetServer("http://open.mapquestapi.com/nominatim/v1")

	// Bring up the http server
	log.Println("\tStarting up the HTTP-Server")
	server := &http.Server{Addr: ":80", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}

type Empty struct{}

type ViewResource struct{}

func (r ViewResource) Register(wsContainer *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/")
	ws.Route(ws.GET("/").To(r.viewHandler))
	wsContainer.Add(ws)
}

func (r ViewResource) viewHandler(req *restful.Request, resp *restful.Response) {
	framecontent, _ := ioutil.ReadFile("index.html")
	stringcontent := string(framecontent)
	resp.ResponseWriter.Write([]byte(stringcontent))
}

type StaticResource struct{}

func (r StaticResource) Register(wsContainer *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/static")
	ws.Route(ws.GET("/{resource:*}").To(r.serveStatic))
	wsContainer.Add(ws)
}

func (r StaticResource) serveStatic(req *restful.Request, resp *restful.Response) {
	path := req.Request.URL.Path
	if path == "/static/js/lang.js" {
		switch req.Request.Header["Accept-Language"][0][0:2] {
		case "de":
			http.ServeFile(resp.ResponseWriter, req.Request, "static/js/translations/de_DE.js")
		default:
			http.ServeFile(resp.ResponseWriter, req.Request, "static/js/translations/en_US.js")
		}
	} else {
		http.ServeFile(resp.ResponseWriter, req.Request, path[1:])
	}
}
