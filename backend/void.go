package main

import (
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"log"
	"net/http"
    "labix.org/v2/mgo"
)

var mongo *mgo.Session

func main() {
	log.Println("Hello world!")
    mng, err := mgo.Dial("localhost")
    mongo = mng
    if err != nil {
        log.Fatal("Could not connect to mongodb!")
    }
	// Prepare REST-backend
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	BuildingResource{}.Register(wsContainer)
	UserResource{}.Register(wsContainer)
	CommentResource{}.Register(wsContainer)
	InvitationResource{}.Register(wsContainer)
	ViewResource{}.Register(wsContainer)
	StaticResource{}.Register(wsContainer)

	// Bring up the http server
	server := &http.Server{Addr: ":80", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}

type ViewResource struct{}

func (r ViewResource) Register(wsContainer *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/")
	ws.Route(ws.GET("/").To(r.viewHandler))
	wsContainer.Add(ws)
}

func (r ViewResource) viewHandler(req *restful.Request, resp *restful.Response) {
	framecontent, _ := ioutil.ReadFile("frame.html")
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
	http.ServeFile(resp.ResponseWriter, req.Request, req.Request.URL.Path[1:])
}
