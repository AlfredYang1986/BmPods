package BmModel
import (

	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type Catenode struct {
	ID           string        `json:"id"`
	Id_          bson.ObjectId `json:"-" bson:"_id"`
	BrandID      string        `json:"brand-id" bson:"brand-id"`
	Status       float64       `json:"status" bson:"status"`
	Prev_cate	 string        `json:"prev-cate" bson:"prev-cate"`		
	Next_cate 	 string        `json:"next-cate" bson:"next-cate"`
	Value		 string        `json:"value" bson:"value"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (c Catenode) GetID() string {
	return c.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (c *Catenode) SetID(id string) error {
	c.ID = id
	return nil
}

func (u *Catenode) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	for k, v := range parameters {
		switch k {
		case "status":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			rst[k] = val
		}
	}
	return rst
}