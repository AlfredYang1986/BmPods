package BmModel

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// Apply is a generic database Apply
type Apply struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	Status       float64 `json:"status" bson:"status"` //0=未处理，1=已处理
	ApplyTime    float64 `json:"apply-time" bson:"apply-time"`
	ExceptTime   float64 `json:"except-time" bson:"except-time"`
	CreateTime   float64 `json:"create-time" bson:"create-time"`
	BrandId      string  `json:"brand-id" bson:"brand-id"`
	ApplyFrom    string  `json:"apply-from" bson:"apply-from"`
	CourseType   float64 `json:"course-type" bson:"course-type"` //0活动 1体验课 2普通课程 -1预注册
	CourseName   string  `json:"course-name" bson:"course-name"`
	Contact      string  `json:"contact" bson:"contact"`
	ReservableId string  `json:"reservable-id" bson:"reservable-id"`

	Kids    []*Kid   `json:"-"`
	KidsIDs []string `json:"kid-ids" bson:"kid-ids"`

	ApplicantID string    `json:"applicant-id" bson:"applicant-id"`
	Applicant   Applicant `json:"-"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Apply) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Apply) SetID(id string) error {
	u.ID = id
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Apply) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type:         "applicants",
			Name:         "applicant",
			Relationship: jsonapi.ToOneRelationship,
		},
		{
			Type:         "kids",
			Name:         "kids",
			Relationship: jsonapi.ToManyRelationship,
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Apply) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, kID := range u.KidsIDs {
		result = append(result, jsonapi.ReferenceID{
			ID:   kID,
			Type: "kids",
			Name: "kids",
		})
	}

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
func (u Apply) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	for key := range u.Kids {
		result = append(result, u.Kids[key])
	}

	if u.ApplicantID != "" {
		result = append(result, u.Applicant)
	}

	return result
}

func (u *Apply) SetToOneReferenceID(name, ID string) error {
	if name == "applicant" {
		u.ApplicantID = ID
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}

// SetToManyReferenceIDs sets the leafs reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (u *Apply) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "kids" {
		u.KidsIDs = IDs
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs adds some new leafs that a users loves so much
func (u *Apply) AddToManyIDs(name string, IDs []string) error {
	if name == "kids" {
		u.KidsIDs = append(u.KidsIDs, IDs...)
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// DeleteToManyIDs removes some leafs from a users because they made him very sick
func (u *Apply) DeleteToManyIDs(name string, IDs []string) error {
	if name == "kids" {
		for _, ID := range IDs {
			for pos, oldID := range u.KidsIDs {
				if ID == oldID {
					// match, this ID must be removed
					u.KidsIDs = append(u.KidsIDs[:pos], u.KidsIDs[pos+1:]...)
				}
			}
		}
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

func (u *Apply) GetConditionsBsonM(parameters map[string][]string) bson.M {

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
		case "course-type":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			rst[k] = val
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
		case "gt[apply-time]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$gt"] = val
			rst["apply-time"] = r
		case "gte[apply-time]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$gte"] = val
			rst["apply-time"] = r
		case "ne[course-type]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$ne"] = val
			rst["course-type"] = r
		}
	}

	return rst
}
