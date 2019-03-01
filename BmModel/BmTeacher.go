package BmModel

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type Teacher struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	Intro       string  `json:"intro" bson:"intro"`
	BrandId     string  `json:"brand-id" bson:"brand-id"`
	Name        string  `json:"name" bson:"name"`
	Nickname    string  `json:"nickname" bson:"nickname"`
	Icon        string  `json:"icon" bson:"icon"`
	Dob         float64 `json:"dob" bson:"dob"`
	Gender      float64 `json:"gender" bson:"gender"`
	RegDate     float64 `json:"reg-date" bson:"reg-date"`
	Contact     string  `json:"contact" bson:"contact"`
	WeChat      string  `json:"wechat" bson:"wechat"`
	JobTitle    string  `json:"job-title" bson:"job-title"`
	JobType     float64 `json:"job-type" bson:"job-type"` //0-兼职, 1-全职
	IdCardNo    string  `json:"id-card-no" bson:"id-card-no"`
	Major       string  `json:"major" bson:"major"`
	TeachYears  float64 `json:"teach-years" bson:"teach-years"`
	Province    string  `json:"province" bson:"province"`
	City        string  `json:"city" bson:"city"`
	District    string  `json:"district" bson:"district"`
	Address     string  `json:"address" bson:"address"`
	NativePlace string  `json:"native-place" bson:"native-place"`
	CreateTime  float64 `bson:"create-time"`
	Archive     float64 `json:"archive" bson:"archive"` //表示在职或离职=归档？
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
	rst["archive"] = float64(0) //不传archive默认只查询存在的，传0只查存在的，传1只查归档的，传-1查全部【包含所有】
	for k, v := range parameters {
		switch k {
		case "brand-id":
			rst[k] = v[0]
		case "archive":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			if val == -1 {
				delete(rst, k)
			} else {
				rst[k] = val
			}
		}
	}
	return rst
}
