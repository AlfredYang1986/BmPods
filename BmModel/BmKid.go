package BmModel

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
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
	Applicant   *Applicant `json:"-"`

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

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Kid) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type:         "applicants",
			Name:         "applicant",
			Relationship: jsonapi.ToOneRelationship,
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Kid) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	if u.ApplicantID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.ApplicantID,
			Type: "applicants",
			Name: "applicant",
		})
	}

	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (u *Kid) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	if u.ApplicantID != "" &&u.Applicant!=nil{
		result = append(result, u.Applicant)
	}

	return result
}

func (u *Kid) SetToOneReferenceID(name, ID string) error {
	if name == "applicant" {
		u.ApplicantID = ID
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}

func (u *Kid) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	r:=make(map[string]interface{})
	var ids []bson.ObjectId	
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
		}
	}
	return rst
}
