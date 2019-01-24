package BmModel

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// Class is a generic database Class
type Class struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	ClassTitle        string  `json:"class-title" bson:"class-title"`
	Status            float64 `json:"status" bson:"status"` //0活动 1体验课 2普通课程
	StartDate         float64 `json:"start-date" bson:"start-date"`
	EndDate           float64 `json:"end-date" bson:"end-date"`
	CreateTime        float64 `json:"create-time" bson:"create-time"`
	CourseTotalCount  float64 `json:"course-total-count"`
	CourseExpireCount float64 `json:"course-expire-count"`
	BrandId           string  `json:"brand-id" bson:"brand-id"`

	Students    []*Student `json:"-"`
	StudentsIDs []string   `json:"-" bson:"student-ids"`
	Teachers    []*Teacher `json:"-"`
	TeachersIDs []string   `json:"-" bson:"teacher-ids"`
	Units       []*Unit    `json:"-"`
	UnitsIDs    []string   `json:"-" bson:"unit-ids"`

	YardID        string      `json:"yard-id" bson:"yard-id"`
	Yard          Yard        `json:"-"`
	SessioninfoID string      `json:"sessioninfo-id" bson:"sessioninfo-id"`
	Sessioninfo   Sessioninfo `json:"-"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Class) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Class) SetID(id string) error {
	u.ID = id
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Class) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "yards",
			Name: "yard",
		},
		{
			Type: "sessioninfos",
			Name: "sessioninfo",
		},
		{
			Type: "students",
			Name: "students",
		},
		{
			Type: "teachers",
			Name: "teachers",
		},
		{
			Type: "units",
			Name: "units",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Class) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, tmpID := range u.StudentsIDs {
		result = append(result, jsonapi.ReferenceID{
			ID:   tmpID,
			Type: "students",
			Name: "students",
		})
	}
	for _, tmpID := range u.TeachersIDs {
		result = append(result, jsonapi.ReferenceID{
			ID:   tmpID,
			Type: "teachers",
			Name: "teachers",
		})
	}
	for _, tmpID := range u.UnitsIDs {
		result = append(result, jsonapi.ReferenceID{
			ID:   tmpID,
			Type: "uints",
			Name: "uints",
		})
	}

	if u.YardID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.YardID,
			Type: "yards",
			Name: "yard",
			Relationship:jsonapi.ToOneRelationship,
		})
	}
	if u.SessioninfoID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.SessioninfoID,
			Type: "sessioninfos",
			Name: "sessioninfo",
		})
	}

	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (u Class) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	for key := range u.Students {
		result = append(result, u.Students[key])
	}
	for key := range u.Teachers {
		result = append(result, u.Teachers[key])
	}
	for key := range u.Units {
		result = append(result, u.Units[key])
	}

	if u.YardID != "" {
		result = append(result, u.Yard)
	}
	if u.SessioninfoID != "" {
		result = append(result, u.Sessioninfo)
	}

	return result
}

func (u *Class) SetToOneReferenceID(name, ID string) error {
	if name == "yard" {
		u.YardID = ID
		return nil
	}
	if name == "sessioninfo" {
		u.SessioninfoID = ID
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}

// SetToManyReferenceIDs sets the leafs reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (u *Class) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "students" {
		u.StudentsIDs = IDs
		return nil
	}
	if name == "teachers" {
		u.TeachersIDs = IDs
		return nil
	}
	if name == "uints" {
		u.UnitsIDs = IDs
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs adds some new leafs that a users loves so much
func (u *Class) AddToManyIDs(name string, IDs []string) error {
	if name == "students" {
		u.StudentsIDs = append(u.StudentsIDs, IDs...)
		return nil
	}
	if name == "teachers" {
		u.TeachersIDs = append(u.TeachersIDs, IDs...)
		return nil
	}
	if name == "uints" {
		u.UnitsIDs = append(u.UnitsIDs, IDs...)
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// DeleteToManyIDs removes some leafs from a users because they made him very sick
func (u *Class) DeleteToManyIDs(name string, IDs []string) error {
	if name == "students" {
		for _, ID := range IDs {
			for pos, oldID := range u.StudentsIDs {
				if ID == oldID {
					// match, this ID must be removed
					u.StudentsIDs = append(u.StudentsIDs[:pos], u.StudentsIDs[pos+1:]...)
				}
			}
		}
	}
	if name == "teachers" {
		for _, ID := range IDs {
			for pos, oldID := range u.TeachersIDs {
				if ID == oldID {
					// match, this ID must be removed
					u.TeachersIDs = append(u.TeachersIDs[:pos], u.TeachersIDs[pos+1:]...)
				}
			}
		}
	}
	if name == "uints" {
		for _, ID := range IDs {
			for pos, oldID := range u.UnitsIDs {
				if ID == oldID {
					// match, this ID must be removed
					u.UnitsIDs = append(u.UnitsIDs[:pos], u.UnitsIDs[pos+1:]...)
				}
			}
		}
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

func (u *Class) GetConditionsBsonM(parameters map[string][]string) bson.M {
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
