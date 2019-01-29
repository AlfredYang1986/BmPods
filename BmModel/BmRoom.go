package BmModel

import "gopkg.in/mgo.v2/bson"

// Room
type Room struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	BrandId  string  `json:"brand-id" bson:"brand-id"`
	Title    string  `json:"title" bson:"title"`
	RoomType float64 `json:"room-type" bson:"room-type"`
	Capacity float64 `json:"capacity" bson:"capacity"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (r Room) GetID() string {
	return r.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (r *Room) SetID(id string) error {
	r.ID = id
	return nil
}

func (u *Room) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	for k, v := range parameters {
		switch k {
		case "brand-id":
			rst[k] = v[0]
		}
	}
	return rst
}
