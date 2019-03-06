package BmModel

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// Reservableitem is a generic database Reservableitem
type Reservableitem struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	BrandId    string  `json:"brand-id" bson:"brand-id"`
	Status     float64 `json:"status" bson:"status"` //0活动 1体验课 2普通课程
	StartDate  float64 `json:"start-date" bson:"start-date"`
	EndDate    float64 `json:"end-date" bson:"end-date"`
	CreateTime float64 `json:"create-time" bson:"create-time"`

	SessioninfoID string      `json:"sessioninfo-id" bson:"sessioninfo-id"`
	Sessioninfo   *Sessioninfo `json:"-"`
	Archive  float64 `json:"archive" bson:"archive"` //表示未结束或已结束=归档？ 
	Execute    float64  //表示是否执行
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Reservableitem) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Reservableitem) SetID(id string) error {
	u.ID = id
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Reservableitem) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "sessioninfos",
			Name: "sessioninfo",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Reservableitem) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	if u.SessioninfoID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.SessioninfoID,
			Type: "sessioninfos",
			Name: "sessioninfo",
		})
	}

	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (u Reservableitem) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	if u.SessioninfoID != ""&&u.Sessioninfo!=nil {
		result = append(result, u.Sessioninfo)
	}

	return result
}

func (u *Reservableitem) SetToOneReferenceID(name, ID string) error {
	if name == "sessioninfo" {
		u.SessioninfoID = ID
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}

func (u *Reservableitem) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	rst["archive"] = float64(0) //不传archive默认只查询存在的，传0只查存在的，传1只查归档的，传-1查全部【包含所有】
	for k, v := range parameters {
		switch k {
		case "brand-id":
			rst[k] = v[0]
		case "status":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			rst[k] = val
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
