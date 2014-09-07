package main

import (
	"github.com/emicklei/go-restful"
	"labix.org/v2/mgo/bson"
    "net/http"
  "errors"
  "math/rand"
  "encoding/hex"
  "encoding/binary"
  "crypto/sha512"
)

type InvitationWrapper struct {
  Invitation *Invitation
}

type Invitation struct {
  MongoId  bson.ObjectId `bson:"_id,omitempty"`
	Id    string `json:"id,omitempty"`
	Email string
  Username string
  Password string
  Organization string
}

func LoadInvitationByEmail(eml string) (*Invitation, error) {
    x := new(Invitation)
    err := mongo.DB("void").C("invitations").Find(bson.M{"email":eml}).One(&x)
    return x,err
}

func LoadInvitationById(id string) (*Invitation, error){
    x := new(Invitation)
    err := mongo.DB("void").C("invitations").Find(bson.M{"id":id}).One(&x)
    return x,err
}

func (i *Invitation) generateId() string {
    s := make([]byte, 8)
    binary.LittleEndian.PutUint64(s,uint64(rand.Int63()))
    s = append(s, []byte(i.MongoId.Hex())...)
    s = append(s, []byte(i.Email)...)
    sum := sha512.Sum512(s)
    return hex.EncodeToString(sum[0:64])
}

func (i *Invitation) Invite () error {
    if _, err := LoadInvitationByEmail(i.Email) ; err == nil {
        return errors.New("invitation to this email already exists")
    }
    i.MongoId = bson.NewObjectId()
    i.Id = i.generateId()
    err := mongo.DB("void").C("invitations").Insert(i)
    return err
}

func (i *Invitation) Activate (d *Invitation) error {
    user := new(User)
    user.Name = d.Username
    user.Email = i.Email
    user.Organization = d.Organization
    user.SetPassword(d.Password)
    if err := user.Save() ; err == nil {
        return mongo.DB("void").C("invitations").RemoveId(i.Id)
    } else {
        return err
    }
}

type InvitationResource struct{}

func (r InvitationResource) Register(wsContainer *restful.Container) {
	ws := new(restful.WebService)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	ws.Path("/rest/invitations")

	ws.Route(ws.GET("/{entry}").Filter(authFilter).To(r.getInvitation))
	ws.Route(ws.POST("/").Filter(authFilter).To(r.createInvitation))
  ws.Route(ws.PUT("/{entry}").Filter(authFilter).To(r.editInvitation))
  ws.Route(ws.DELETE("/{entry}").Filter(authFilter).To(r.deleteInvitation))
	wsContainer.Add(ws)
}

func (r InvitationResource) getInvitation(req *restful.Request, resp *restful.Response) {
  if inv, err := LoadInvitationById(req.PathParameter("entry")); err == nil {
      iw := new(InvitationWrapper)
      iw.Invitation = inv
      resp.WriteEntity(iw)
  } else {
      resp.WriteErrorString(http.StatusNotFound, "no such invitation")
  }
}

func (r InvitationResource) createInvitation(req *restful.Request, resp *restful.Response) {
    reqUser := getRequestUser(req)
    if reqUser == nil {
        resp.WriteErrorString(http.StatusForbidden, "you must be logged in to do that")
        return
    }
    iw := new(InvitationWrapper)
    if err := req.ReadEntity(iw) ; err == nil {
        if err := iw.Invitation.Invite() ; err != nil {
            resp.WriteErrorString(http.StatusInternalServerError, err.Error())
        } else {
            resp.WriteEntity(&iw)
        }
    }else{
        resp.WriteErrorString(http.StatusBadRequest, "bad invitation")
    }
}

func (r InvitationResource) editInvitation (req *restful.Request, resp *restful.Response) {
  if inv, err := LoadInvitationById(req.PathParameter("entry")); err == nil {
      iw := new(InvitationWrapper)
      if err := req.ReadEntity(iw) ; err == nil {
          if err := inv.Activate(iw.Invitation) ; err != nil{
              resp.WriteErrorString(http.StatusBadRequest, err.Error())
          }
      } else {
        resp.WriteErrorString(http.StatusBadRequest, "bad invitation")
      }
  } else {
      resp.WriteErrorString(http.StatusNotFound, "no such invitation")
  }
}

func (r InvitationResource) deleteInvitation (req *restful.Request, resp *restful.Response) {
    resp.WriteEntity(Empty{})
}
