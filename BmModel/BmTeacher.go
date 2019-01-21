package BmModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Teacher struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	Intro string `json:"intro" bson:"intro"`

	BrandId string `json:"brand-id" bson:"brand-id"`

	Name       string  `json:"name" bson:"name"`
	Nickname   string  `json:"nickname" bson:"nickname"`
	Icon       string  `json:"icon" bson:"icon"`
	Dob        float64 `json:"dob" bson:"dob"`
	Gender     float64 `json:"gender" bson:"gender"`
	RegDate    float64 `json:"reg-date" bson:"reg-date"`
	Contact    string  `json:"contact" bson:"contact"`
	WeChat     string  `json:"wechat" bson:"wechat"`
	JobTitle   string  `json:"job-title" bson:"job-title"`
	JobType    float64 `json:"job-type" bson:"job-type"` //0-兼职, 1-全职
	IdCardNo   string  `json:"id-card-no" bson:"id-card-no"`
	Major      string  `json:"major" bson:"major"`
	TeachYears float64 `json:"teach-years" bson:"teach-years"`

	Province    string `json:"province" bson:"province"`
	City        string `json:"city" bson:"city"`
	District    string `json:"district" bson:"district"`
	Address     string `json:"address" bson:"address"`
	NativePlace string `json:"native-place" bson:"native-place"`

	CreateTime float64 `bson:"create-time"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (c Teacher) GetID() string {
	return c.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (c *Teacher) SetID(id string) error {
	c.ID = id
	return nil
}

func (u *Teacher) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	for k, v := range parameters {
		switch k {
		case "brand-id":
			rst[k] = v[0]
		}
	}
	return rst
}
