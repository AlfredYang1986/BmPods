package BmModel

import (
	"gopkg.in/mgo.v2/bson"
)

type ClassUnitBind struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	UnitID    string  `json:"unit-id" bson:"unit-id"`
	ClassID   string  `json:"class-id" bson:"class-id"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u ClassUnitBind) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *ClassUnitBind) SetID(id string) error {
	u.ID = id
	return nil
}

func (u *ClassUnitBind) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	for k, v := range parameters {
		switch k {
		case "class-id":
			rst[k] = v[0]
		}
	}
	return rst
}
