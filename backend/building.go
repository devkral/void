package main

import (
	"github.com/emicklei/go-restful"
	"github.com/grindhold/gominatim"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
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

type BuildingsWrapper struct {
	Buildings []*Building
}

type BuildingWrapper struct {
	Building *Building
}

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

	Comments []bson.ObjectId `json:"comments"`
}

func LoadBuildingById(id bson.ObjectId) (*Building, error) {
	x := new(Building)
	err := mongo.DB("void").C("buildings").Find(bson.M{"_id": id}).One(x)
	return x, err
}

func LoadBuildings() ([]*Building, error) {
	x := make([]*Building, 0)
	err := mongo.DB("void").C("buildings").Find(bson.M{}).All(&x)
	return x, err
}

func (b *Building) Save() error {
	if !b.Id.Valid() {
		b.Id = bson.NewObjectId()
		qry := new(gominatim.SearchQuery)
		qry.Q = b.Street + " " + b.Number + ", " + b.Zip + " " + b.City
		res, err := qry.Get()
		if err == nil {
			b.Lat = res[0].Lat
			b.Lat_f, _ = strconv.ParseFloat(b.Lat, 64)
			b.Lon = res[0].Lon
			b.Lon_f, _ = strconv.ParseFloat(b.Lon, 64)
		} else {
			log.Println("Err: " + err.Error())
		}
	} else {
	}
	_, err := mongo.DB("void").C("buildings").UpsertId(b.Id, b)
	return err
}

func (b *Building) Update(u *Building) {
	if b.Street != u.Street || b.Number != u.Number || b.City != u.City || b.Zip != u.Zip {
		qry := new(gominatim.SearchQuery)
		qry.Q = b.Street + " " + b.Number + ", " + b.Zip + " " + b.City
		res, err := qry.Get()
		if err == nil {
			b.Lat = res[0].Lat
			b.Lat_f, _ = strconv.ParseFloat(b.Lat, 64)
			b.Lon = res[0].Lon
			b.Lon_f, _ = strconv.ParseFloat(b.Lon, 64)
		}
	}
	b.Street = u.Street
	b.Number = u.Number
	b.City = u.City
	b.Zip = u.Zip
	b.OwnerName = u.OwnerName
	b.OwnerPhone = u.OwnerPhone
	b.OwnerEmail = u.OwnerEmail

	b.Area = u.Area

	b.Description = u.Description

	b.Status = u.Status
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
	ws.Route(ws.GET("/{entry}").Filter(authFilter).To(r.getBuilding))
	ws.Route(ws.POST("/").Filter(authFilter).To(r.createBuilding))
	ws.Route(ws.PUT("/{entry}").Filter(authFilter).To(r.editBuilding))
	wsContainer.Add(ws)
}

func (r BuildingResource) getBuildings(req *restful.Request, resp *restful.Response) {
    reqUser := getRequestUser(req)
    if reqUser == nil {
        resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
        return
    }
	if buildings, err := LoadBuildings(); err == nil {
		bw := new(BuildingsWrapper)
		bw.Buildings = buildings
		resp.WriteEntity(bw)
	} else {
		resp.WriteErrorString(http.StatusInternalServerError, "Nothing Found")
	}
}

func (r BuildingResource) getBuilding(req *restful.Request, resp *restful.Response) {
    reqUser := getRequestUser(req)
    if reqUser == nil {
        resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
        return
    }
	b, err := LoadBuildingById(bson.ObjectIdHex(req.PathParameter("entry")))
	if err != nil {
		resp.WriteErrorString(http.StatusNotFound, "no such building")
	} else {
		bw := new(BuildingWrapper)
		bw.Building = b
		resp.WriteEntity(bw)
	}
}

func (r BuildingResource) createBuilding(req *restful.Request, resp *restful.Response) {
	bw := new(BuildingWrapper)
	err := req.ReadEntity(bw)
	if err == nil {
		bw.Building.Save()
		resp.WriteEntity(bw)
	} else {
		resp.WriteErrorString(http.StatusBadRequest, "Your building is invalid")
	}
}

func (r BuildingResource) editBuilding(req *restful.Request, resp *restful.Response) {
    reqUser := getRequestUser(req)
    if reqUser == nil {
        resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
        return
    }
	bw := new(BuildingWrapper)
	err := req.ReadEntity(bw)
	if err == nil {
		b, err := LoadBuildingById(bson.ObjectIdHex(req.PathParameter("entry")))
		if err != nil {
			resp.WriteErrorString(http.StatusNotFound, "Cannot edit nonexistent building")
			return
		}
		b.Update(bw.Building)
		b.Save()
		resp.WriteEntity(bw)
	} else {
		resp.WriteErrorString(http.StatusBadRequest, "Your building is invalid")
	}
}
