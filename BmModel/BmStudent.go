package BmModel

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/manyminds/api2go/jsonapi"
	"errors"
)

type Student struct {
	ID  string        `json:"id"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	StudentId     string `json:"brand-id" bson:"brand-id"`

	// 顾问， 关联关系 teacher
	//TeacherId   string `json:"teacher-id" bson:"teacher-id"`
	//TeacherName string `json:"teacher-name" bson:"teacher-name"`

	SourceWay   string `json:"source-way" bson:"source-way"`	// 来源于

	Intro       string  `json:"intro" bson:"intro"`
	Status      string  `json:"status" bson:"status"` //未付款-candidate, 已付款-stud
	LessonCount float64 `json:"lesson-count" bson:"lesson-count"`

	Name        string  `json:"name" bson:"name"`
	Nickname    string  `json:"nickname" bson:"nickname"`
	Icon        string  `json:"icon" bson:"icon"`
	Dob         float64 `json:"dob" bson:"dob"`
	Gender      float64 `json:"gender" bson:"gender"`
	RegDate     float64 `json:"reg-date" bson:"reg-date"`
	CreateTime  float64 `json:"create-time" bson:"create-time"`
	Contact     string  `json:"contact" bson:"contact"`
	WeChat      string  `json:"wechat" bson:"wechat"`

	Province    string  `json:"province" bson:"province"`
	City        string  `json:"city" bson:"city"`
	District    string  `json:"district" bson:"district"`
	Address     string  `json:"address" bson:"address"`
	School      string  `json:"school" bson:"school"`
	IdCardNo    string  `json:"id-card-no" bson:"id-card-no"`

	KidID string `json:"-" bson:"kid-id"`
	Kid *Kid `json:"-"`

	Teacher *Teacher `json:"-"`
	TeacherID string `json:"-" bson:"teacher-id"`

	Guardians []*Guardian `json:"-"`
	GuardiansIDs []string `json:"-" bson:"guardians-ids"`
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
			Type: "Kid",
			Name: "kid",
		},
		{
			Type: "Teacher",
			Name: "teacher",
		},
		{
			Type: "Guardian",
			Name: "guardians",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Student) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, kID := range u.GuardiansIDs {
		result = append(result, jsonapi.ReferenceID{
			ID:   kID,
			Type: "Guardian",
			Name: "guardian",
		})
	}

	if u.TeacherID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.TeacherID,
			Type: "Teacher",
			Name: "teacher",
		})
	}

	if u.KidID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID: u.KidID,
			Type: "Kid",
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

	if u.TeacherID != "" {
		result = append(result, u.Teacher)
	}

	if u.KidID != "" {
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
