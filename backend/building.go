package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo/bson"
)

// Possible states of a building

const (
	S_SUSPECTED = iota
	S_CONFIRMED
	S_SOLD
	S_TORENT
	S_RENTED
	S_TOBUY
	S_OCCUPIED
	S_PACHT    //TODO: translate to english
	S_ERBPACHT //TODO: translate to english
)

type Building struct {
	Id     bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Street string
	Number string
	City   string
	Zip    string

	Lat   string
	Lat_f float64
	Lon   string
	Lon_f float64

	OwnerName  string
	OwnerPhone string
	OwnerEmail string

	Area int

	Description string

	Status int

	Comments []bson.ObjectId
}

func LoadBuildingById(id *bson.ObjectId) (*Building, error) {
	x := new(Building)
	err := mongo.DB("void").C("buildings").Find(bson.M{"_id": id}).One(x)
	return x, err
}

func LoadBuildings() ([]*Building, error) {
	x := make([]*Building, 0)
	err := mongo.DB("void").C("buildings").Find(bson.M{}).All(x)
	return x, err
}

func (b *Building) Save() error {
	if !b.Id.Valid() {
		b.Id = bson.NewObjectId()
	} else {
	}
	_, err := mongo.DB("void").C("buildings").UpsertId(b.Id, b)
	return err
}

func (b *Building) Update(update *Building) error {
	return nil
}

func (b *Building) Delete() error {
	return mongo.DB("void").C("buildings").RemoveId(b.Id)
}

type BuildingResource struct{}

func (r BuildingResource) Register(wsContainer *restful.Container) {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path("/rest/buildings")

	ws.Route(ws.GET("/").Filter(authFilter).To(r.getBuildings))
	ws.Route(ws.POST("/").Filter(authFilter).To(r.createBuilding))
	ws.Route(ws.PUT("/{entry}").Filter(authFilter).To(r.editBuilding))
    wsContainer.Add(ws)
}

func (r BuildingResource) getBuildings(req *restful.Request, resp *restful.Response) {

}

func (r BuildingResource) createBuilding(req *restful.Request, resp *restful.Response) {

}

func (r BuildingResource) editBuilding(req *restful.Request, resp *restful.Response) {

}
