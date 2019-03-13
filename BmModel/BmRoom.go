package BmModel

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// Room
type Room struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	BrandId  string  `json:"brand-id" bson:"brand-id"`
	Title    string  `json:"title" bson:"title"`
	RoomType float64 `json:"room-type" bson:"room-type"`
	Capacity float64 `json:"capacity" bson:"capacity"`
	Archive  float64 `json:"archive" bson:"archive"` //表示是否归档？ 
	IsUnit   float64 `json:"isunit" bson:"isunit"`   //表示未排课或已排课=归档？ 
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
	rst["archive"] = float64(0) //不传archive默认只查询开放的，传0只查开放的，传1只查归档的，传-1查全部【包含所有】
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
		case "roomids":
			r:=make(map[string]interface{})
			var ids []bson.ObjectId
			for i:=0;i<len(v);i++{
				ids=append(ids,bson.ObjectIdHex(v[i]))
			}
			r["$in"]=ids
			rst["_id"] = r
		}
	}
	return rst
}
