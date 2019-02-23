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
	Sessioninfo   Sessioninfo `json:"-"`

	Classes    []*Class `json:"-"`
	ClassesIDs []string `json:"-" bson:"class-ids"`
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
		{
			Type: "classes",
			Name: "classes",
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

	for _, tmpID := range u.ClassesIDs {
		result = append(result, jsonapi.ReferenceID{
			ID:   tmpID,
			Type: "classes",
			Name: "classes",
		})
	}

	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (u Reservableitem) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	if u.SessioninfoID != "" {
		result = append(result, u.Sessioninfo)
	}
	for key := range u.Classes {
		result = append(result, u.Classes[key])
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

// SetToManyReferenceIDs sets the leafs reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (u *Reservableitem) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "classes" {
		u.ClassesIDs = IDs
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs adds some new leafs that a users loves so much
func (u *Reservableitem) AddToManyIDs(name string, IDs []string) error {
	if name == "classes" {
		u.ClassesIDs = append(u.ClassesIDs, IDs...)
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// DeleteToManyIDs removes some leafs from a users because they made him very sick
func (u *Reservableitem) DeleteToManyIDs(name string, IDs []string) error {
	if name == "classes" {
		for _, ID := range IDs {
			for pos, oldID := range u.ClassesIDs {
				if ID == oldID {
					// match, this ID must be removed
					u.ClassesIDs = append(u.ClassesIDs[:pos], u.ClassesIDs[pos+1:]...)
				}
			}
		}
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

func (u *Reservableitem) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
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
		}
	}
	return rst
}
