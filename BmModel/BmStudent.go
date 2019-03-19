package BmModel

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type Student struct {
	ID  string        `json:"id"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	BrandId string `json:"brand-id" bson:"brand-id"`

	SourceWay   string  `json:"source-way" bson:"source-way"` // 来源于
	Intro       string  `json:"intro" bson:"intro"`
	Status      float64 `json:"status" bson:"status"` //0-潜在, 1-正式, -1-休学
	LessonCount float64 `json:"lesson-count" bson:"lesson-count"`

	Name       string  `json:"name" bson:"name"`
	Nickname   string  `json:"nickname" bson:"nickname"`
	Icon       string  `json:"icon" bson:"icon"`
	Dob        float64 `json:"dob" bson:"dob"`
	Gender     float64 `json:"gender" bson:"gender"`
	RegDate    float64 `json:"reg-date" bson:"reg-date"`
	CreateTime float64 `json:"create-time" bson:"create-time"`
	Contact    string  `json:"contact" bson:"contact"`
	WeChat     string  `json:"wechat" bson:"wechat"`

	Province string `json:"province" bson:"province"`
	City     string `json:"city" bson:"city"`
	District string `json:"district" bson:"district"`
	Address  string `json:"address" bson:"address"`
	School   string `json:"school" bson:"school"`
	Grade    string `json:"grade" bson:"grade"`
	IdCardNo string `json:"id-card-no" bson:"id-card-no"`

	KidID string `json:"kid-id" bson:"kid-id"`
	Kid   *Kid   `json:"-"`

	Teacher   *Teacher `json:"-"`
	TeacherID string   `json:"-" bson:"teacher-id"`

	Guardians    []*Guardian `json:"-"`
	GuardiansIDs []string    `json:"-" bson:"guardian-ids"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Student) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Student) SetID(id string) error {
	u.ID = id
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Student) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "kids",
			Name: "kid",
		},
		{
			Type: "teachers",
			Name: "teacher",
		},
		{
			Type: "guardians",
			Name: "guardians",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Student) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, gID := range u.GuardiansIDs {
		result = append(result, jsonapi.ReferenceID{
			ID:   gID,
			Type: "guardians",
			Name: "guardians",
		})
	}

	if u.TeacherID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.TeacherID,
			Type: "teachers",
			Name: "teacher",
		})
	}

	if u.KidID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.KidID,
			Type: "kids",
			Name: "kid",
		})
	}

	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (u Student) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	for key := range u.Guardians {
		result = append(result, u.Guardians[key])
	}

	if u.TeacherID != "" && u.Teacher != nil {
		result = append(result, u.Teacher)
	}

	if u.KidID != "" && u.Kid != nil {
		result = append(result, u.Kid)
	}
	return result
}

func (u *Student) SetToOneReferenceID(name, ID string) error {
	if name == "kid" {
		u.KidID = ID
		return nil
	}

	if name == "teacher" {
		u.TeacherID = ID
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}

// SetToManyReferenceIDs sets the leafs reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (u *Student) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "guardians" {
		u.GuardiansIDs = IDs
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs adds some new leafs that a users loves so much
func (u *Student) AddToManyIDs(name string, IDs []string) error {
	if name == "guardians" {
		u.GuardiansIDs = append(u.GuardiansIDs, IDs...)
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// DeleteToManyIDs removes some leafs from a users because they made him very sick
func (u *Student) DeleteToManyIDs(name string, IDs []string) error {
	if name == "guardians" {
		for _, ID := range IDs {
			for pos, oldID := range u.GuardiansIDs {
				if ID == oldID {
					// match, this ID must be removed
					u.GuardiansIDs = append(u.GuardiansIDs[:pos], u.GuardiansIDs[pos+1:]...)
				}
			}
		}
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

func (u *Student) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	r := make(map[string]interface{})
	var ids []bson.ObjectId
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
		case "studentsids":
			for i := 0; i < len(v); i++ {
				ids = append(ids, bson.ObjectIdHex(v[i]))
			}
			r["$in"] = ids
			rst["_id"] = r
		case "kidids":
			var kidids []string
			for i := 0; i < len(v); i++ {
				kidids = append(kidids, v[i])
			}
			r["$in"] = kidids
			rst["kid-id"] = r
		}
	}
	return rst
}
