package BmModel

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// Kid is the Kid that a user consumes in order to get fat and happy
type Kid struct {
	ID           string        `json:"id"`
	Id_          bson.ObjectId `json:"-" bson:"_id"`
	Name         string        `json:"name" bson:"name"`
	NickName     string        `json:"nickname" bson:"nickname"`
	Gender       float64       `json:"gender" bson:"gender"`
	Dob          float64       `json:"dob" bson:"dob"`
	GuardianRole string        `json:"guardian-role" bson:"guardian-role"`

	ApplicantID string    `json:"applicant-id" bson:"applicant-id"`
	Applicant   Applicant `json:"-"`
	Archive     float64 `json:"archive" bson:"archive"` //表示归档？
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (c Kid) GetID() string {
	return c.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (c *Kid) SetID(id string) error {
	c.ID = id
	return nil
}

func (u *Kid) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	r:=make(map[string]interface{})
	var ids []bson.ObjectId
	rst["archive"] = float64(0) //不传archive默认只查询存在的，传0只查存在的，传1只查归档的，传-1查全部【包含所有】
	for k, v := range parameters {
		switch k {
		case "applicant-id":
			rst[k] = v[0]	
		case "kidsids":
			for i:=0;i<len(v);i++{
				ids=append(ids,bson.ObjectIdHex(v[i]))
			}
			r["$in"]=ids
			rst["_id"] = r
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
