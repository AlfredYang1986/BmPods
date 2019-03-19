package BmModel

import "gopkg.in/mgo.v2/bson"

// Guardian is the Guardian that a user consumes in order to get fat and happy
type Guardian struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	RelationShip string  `json:"relation-ship" bson:"relation-ship"`
	BrandId      string  `json:"brand-id" bson:"brand-id"`
	Name         string  `json:"name" bson:"name"`
	Nickname     string  `json:"nickname" bson:"nickname"`
	Icon         string  `json:"icon" bson:"icon"`
	Dob          float64 `json:"dob" bson:"dob"`
	Gender       float64 `json:"gender" bson:"gender"`
	RegDate      float64 `json:"reg-date" bson:"reg-date"`
	Contact      string  `json:"contact" bson:"contact"`
	WeChat       string  `json:"wechat" bson:"wechat"`
	Address      string  `json:"address" bson:"address"`
	IdCardNo     string  `json:"id-card-no" bson:"id-card-no"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (c Guardian) GetID() string {
	return c.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (c *Guardian) SetID(id string) error {
	c.ID = id
	return nil
}

func (u *Guardian) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	r:=make(map[string]interface{})
	var ids []bson.ObjectId
	for k, v := range parameters {
		switch k {
		case "guardiansids":
			for i:=0;i<len(v);i++{
				ids=append(ids,bson.ObjectIdHex(v[i]))
			}
			r["$in"]=ids
			rst["_id"] = r
		case "contact":
			rst[k] = v[0]
		}
	}
	return rst
}
