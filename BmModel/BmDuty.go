package BmModel

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
)

type Duty struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	TeacherDuty    string  `json:"teacher-duty" bson:"teacher-duty"`

	TeacherID string  `json:"teacher-id" bson:"teacher-id"`
	Teacher   Teacher `json:"-"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Duty) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Duty) SetID(id string) error {
	u.ID = id
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Duty) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "teachers",
			Name: "teacher",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Duty) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	if u.TeacherID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.TeacherID,
			Type: "teachers",
			Name: "teacher",
		})
	}

	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (u Duty) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	if u.TeacherID != "" {
		result = append(result, u.Teacher)
	}

	return result
}

func (u *Duty) SetToOneReferenceID(name, ID string) error {
	if name == "teacher" {
		u.TeacherID = ID
		return nil
	}
	return errors.New("There is no to-one relationship with the name " + name)
}

func (u *Duty) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	return rst
}
