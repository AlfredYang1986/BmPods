package BmModel

import "gopkg.in/mgo.v2/bson"

// Category is the Category that a user consumes in order to get fat and happy
type Category struct {
	ID       string        `json:"-"`
	Id_      bson.ObjectId `json:"-" bson:"_id"`
	Title    string        `json:"title" bson:"title"`
	SubTitle string        `json:"sub-title" bson:"sub-title"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (c Category) GetID() string {
	return c.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (c *Category) SetID(id string) error {
	c.ID = id
	return nil
}

func (u *Category) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return  bson.M {}
}
