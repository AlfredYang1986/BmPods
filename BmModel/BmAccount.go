package BmModel

import (
	"github.com/alfredyang1986/blackmirror/bmmodel"
	"github.com/alfredyang1986/blackmirror/bmmodel/request"
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	Account  string `json:"account" bson:"account"`
	Password string `json:"password" bson:"password"`
	BrandId  string `json:"brand-id" bson:"brand-id"`
	Token    string `json:"token"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Account) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Account) SetID(id string) error {
	u.ID = id
	return nil
}

func (u *Account) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}

// TODO 老Model写法
/*------------------------------------------------
 * bm object interface
 *------------------------------------------------*/

func (bd *Account) ResetIdWithId_() {
	bmmodel.ResetIdWithId_(bd)
}

func (bd *Account) ResetId_WithID() {
	bmmodel.ResetId_WithID(bd)
}

/*------------------------------------------------
 * bmobject interface
 *------------------------------------------------*/

func (bd *Account) QueryObjectId() bson.ObjectId {
	return bd.Id_
}

func (bd *Account) QueryId() string {
	return bd.ID
}

func (bd *Account) SetObjectId(id_ bson.ObjectId) {
	bd.Id_ = id_
}

func (bd *Account) SetId(id string) {
	bd.ID = id
}

/*------------------------------------------------
 * relationships interface
 *------------------------------------------------*/
func (bd Account) SetConnect(tag string, v interface{}) interface{} {
	return bd
}

func (bd Account) QueryConnect(tag string) interface{} {
	return bd
}

/*------------------------------------------------
 * mongo interface
 *------------------------------------------------*/

func (bd *Account) InsertBMObject() error {
	return bmmodel.InsertBMObject(bd)
}

func (bd *Account) FindOne(req request.Request) error {
	return bmmodel.FindOne(req, bd)
}

func (bd *Account) UpdateBMObject(req request.Request) error {
	return bmmodel.UpdateOne(req, bd)
}
