package BmModel

import (
	"github.com/alfredyang1986/blackmirror/bmmodel"
	"github.com/alfredyang1986/blackmirror/bmmodel/request"
	"gopkg.in/mgo.v2/bson"
)

type Applicant struct {
	ID  string        `json:"id"`
	Id_ bson.ObjectId `bson:"_id"`

	Name            string  `json:"name" bson:"name"`
	Gender          float64 `json:"gender" bson:"gender"`
	Pic             string  `json:"pic" bson:"pic"`
	RegisterPhone   string  `json:"regi-phone" bson:"regi-phone"`
	WeChatOpenid    string  `json:"wechat-openid" bson:"wechat-openid"`
	UnionID         string  `json:"union-id" bson:"union-id"`
	WeChatBindPhone string  `json:"wechat-bind-phone" bson:"wechat-bind-phone"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Applicant) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Applicant) SetID(id string) error {
	u.ID = id
	return nil
}

func (u *Applicant) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}

// TODO 老Model写法
/*------------------------------------------------
 * bm object interface
 *------------------------------------------------*/

func (bd *Applicant) ResetIdWithId_() {
	bmmodel.ResetIdWithId_(bd)
}

func (bd *Applicant) ResetId_WithID() {
	bmmodel.ResetId_WithID(bd)
}

/*------------------------------------------------
 * bmobject interface
 *------------------------------------------------*/

func (bd *Applicant) QueryObjectId() bson.ObjectId {
	return bd.Id_
}

func (bd *Applicant) QueryId() string {
	return bd.ID
}

func (bd *Applicant) SetObjectId(id_ bson.ObjectId) {
	bd.Id_ = id_
}

func (bd *Applicant) SetId(id string) {
	bd.ID = id
}

/*------------------------------------------------
 * relationships interface
 *------------------------------------------------*/
func (bd Applicant) SetConnect(tag string, v interface{}) interface{} {
	return bd
}

func (bd Applicant) QueryConnect(tag string) interface{} {
	return bd
}

/*------------------------------------------------
 * mongo interface
 *------------------------------------------------*/

func (bd *Applicant) InsertBMObject() error {
	return bmmodel.InsertBMObject(bd)
}

func (bd *Applicant) FindOne(req request.Request) error {
	return bmmodel.FindOne(req, bd)
}

func (bd *Applicant) UpdateBMObject(req request.Request) error {
	return bmmodel.UpdateOne(req, bd)
}
