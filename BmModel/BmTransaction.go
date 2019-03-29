package BmModel

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// Transaction is a generic database Transaction
type Transaction struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	BrandId         string  `json:"brand-id" bson:"brand-id"`
	OrderId         string  `json:"order-id" bson:"order-id"`
	OrderTime       float64 `json:"order-time" bson:"order-time"`
	CreateTime      float64 `json:"create-time" bson:"create-time"`
	OrderWay        string  `json:"order-way" bson:"order-way"`
	MoneyReceivable float64 `json:"money-receivable" bson:"money-receivable"`
	MoneyReceived   float64 `json:"money-received" bson:"money-received"`
	MoneyUnit       string  `json:"money-unit" bson:"money-unit"`
	Payment         string  `json:"payment" bson:"payment"`
	Remark          string  `json:"remark" bson:"remark"`
	Operator          string  `json:"operator" bson:"operator"`

	Attachables    []*Attachable `json:"-"`
	AttachablesIDs []string      `json:"attachable-ids" bson:"attachable-ids"`
	ApplicantID    string        `json:"applicant-id" bson:"applicant-id"`
	Applicant      *Applicant    `json:"-"`
	TeacherID      string        `json:"teacher-id" bson:"teacher-id"`
	Teacher        *Teacher      `json:"-"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Transaction) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Transaction) SetID(id string) error {
	u.ID = id
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Transaction) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type:         "applicants",
			Name:         "applicant",
			Relationship: jsonapi.ToOneRelationship,
		},
		{
			Type:         "teachers",
			Name:         "teacher",
			Relationship: jsonapi.ToOneRelationship,
		},
		{
			Type:         "attachables",
			Name:         "attachables",
			Relationship: jsonapi.ToManyRelationship,
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Transaction) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, aID := range u.AttachablesIDs {
		result = append(result, jsonapi.ReferenceID{
			ID:   aID,
			Type: "attachables",
			Name: "attachables",
		})
	}

	if u.ApplicantID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.ApplicantID,
			Type: "applicants",
			Name: "applicant",
		})
	}
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
func (u Transaction) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	for key := range u.Attachables {
		result = append(result, u.Attachables[key])
	}

	if u.ApplicantID != "" && u.Applicant != nil {
		result = append(result, u.Applicant)
	}
	if u.TeacherID != "" && u.Teacher != nil {
		result = append(result, u.Teacher)
	}

	return result
}

func (u *Transaction) SetToOneReferenceID(name, ID string) error {
	if name == "applicant" {
		u.ApplicantID = ID
		return nil
	}
	if name == "teacher" {
		u.TeacherID = ID
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}

// SetToManyReferenceIDs sets the leafs reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (u *Transaction) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "attachables" {
		u.AttachablesIDs = IDs
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs adds some new leafs that a users loves so much
func (u *Transaction) AddToManyIDs(name string, IDs []string) error {
	if name == "attachables" {
		u.AttachablesIDs = append(u.AttachablesIDs, IDs...)
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// DeleteToManyIDs removes some leafs from a users because they made him very sick
func (u *Transaction) DeleteToManyIDs(name string, IDs []string) error {
	if name == "attachables" {
		for _, ID := range IDs {
			for pos, oldID := range u.AttachablesIDs {
				if ID == oldID {
					// match, this ID must be removed
					u.AttachablesIDs = append(u.AttachablesIDs[:pos], u.AttachablesIDs[pos+1:]...)
				}
			}
		}
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

func (u *Transaction) GetConditionsBsonM(parameters map[string][]string) bson.M {

	rst := make(map[string]interface{})
	for k, v := range parameters {
		switch k {
		case "brand-id":
			rst[k] = v[0]
		case "applicant-id":
			rst[k] = v[0]
		case "teacher-id":
			rst[k] = v[0]
		case "order-id":
			rst[k] = v[0]
		case "lt[create-time]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$lt"] = val
			rst["create-time"] = r
		case "lte[create-time]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$lte"] = val
			rst["create-time"] = r
		case "gt[order-time]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$gt"] = val
			rst["order-time"] = r
		case "gte[order-time]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$gte"] = val
			rst["order-time"] = r
		}
	}

	return rst
}
