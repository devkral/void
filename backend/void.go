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
	gominatim.SetServer("http://nominatim.openstreetmap.org/")

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
	http.ServeFile(resp.ResponseWriter, req.Request, req.Request.URL.Path[1:])
}
