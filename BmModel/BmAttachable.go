package BmModel

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
)

// Attachable is a generic database Attachable
type Attachable struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	BrandId       string `json:"brand-id" bson:"brand-id"`
	SessioninfoID string `json:"sessioninfo-id" bson:"sessioninfo-id"`
	ReservableID  string `json:"reservable-id" bson:"reservable-id"`

	Title               string  `json:"title" bson:"title"`
	Alb                 float64 `json:"alb" bson:"alb"`
	Aub                 float64 `json:"aub" bson:"aub"`
	StandardCourseCount float64 `json:"standard-course-count" bson:"standard-course-count"`
	StandardPrice       float64 `json:"standard-price" bson:"standard-price"`
	StandardPriceUnit   string  `json:"standard-price-unit" bson:"standard-price-unit"`
	PreferentialPrice   float64 `json:"preferential-price" bson:"preferential-price"`
	SignedPrice         float64 `json:"signed-price" bson:"signed-price"`
	Archive             float64 `json:"archive" bson:"archive"` //表示是否归档

	Student   *Student `json:"-"`
	StudentID string   `json:"-" bson:"student-ids"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Attachable) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Attachable) SetID(id string) error {
	u.ID = id
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Attachable) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "students",
			Name: "student",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Attachable) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	if u.StudentID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.StudentID,
			Type: "students",
			Name: "student",
		})
	}
	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (u Attachable) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	if u.StudentID != "" && u.Student != nil {
		result = append(result, u.Student)
	}

	return result
}

func (u *Attachable) SetToOneReferenceID(name, ID string) error {
	if name == "student" {
		u.StudentID = ID
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}

func (u *Attachable) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	for k, v := range parameters {
		switch k {
		case "attachableids":
			r := make(map[string]interface{})
			var ids []bson.ObjectId
			for i := 0; i < len(v); i++ {
				ids = append(ids, bson.ObjectIdHex(v[i]))
			}
			r["$in"] = ids
			rst["_id"] = r
		}
	}
	return rst
}
