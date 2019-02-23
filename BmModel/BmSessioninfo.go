package BmModel

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

// Sessioninfo is a generic database Sessioninfo
type Sessioninfo struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	BrandId      string   `json:"brand-id" bson:"brand-id"`
	Title        string   `json:"title" bson:"title"`
	SubTitle     string   `json:"sub-title" bson:"sub-title"`
	Alb          float64  `json:"alb" bson:"alb"`
	Aub          float64  `json:"aub" bson:"aub"`
	Level        string   `json:"level" bson:"level"`
	Count        float64  `json:"count" bson:"count"`
	Length       float64  `json:"length" bson:"length"`
	Description  string   `json:"description" bson:"description"`
	Harvest      string   `json:"harvest" bson:"harvest"`
	Accompany    float64  `json:"accompany" bson:"accompany"`
	Status       float64  `json:"status" bson:"status"` //0活动 1体验课 2普通课程
	Acquisition  []string `json:"acquisition" bson:"acquisition"`
	Including    []string `json:"including" bson:"including"`
	Carrying     []string `json:"carrying" bson:"carrying"`
	Notice       string   `json:"notice" bson:"notice"`
	PlayChildren string   `json:"play-children" bson:"play-children"`
	Cover        string   `json:"cover" bson:"cover"`

	Images    []*Image `json:"-"`
	ImagesIDs []string `json:"-" bson:"image-ids"`

	Category   Category `json:"-"`
	CategoryID string   `json:"category-id" bson:"category-id"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Sessioninfo) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Sessioninfo) SetID(id string) error {
	u.ID = id
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Sessioninfo) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "categories",
			Name: "category",
		},
		{
			Type: "images",
			Name: "images",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Sessioninfo) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, kID := range u.ImagesIDs {
		result = append(result, jsonapi.ReferenceID{
			ID:   kID,
			Type: "images",
			Name: "images",
		})
	}

	if u.CategoryID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.CategoryID,
			Type: "categories",
			Name: "category",
		})
	}

	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (u Sessioninfo) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	for key := range u.Images {
		result = append(result, u.Images[key])
	}

	if u.CategoryID != "" {
		result = append(result, u.Category)
	}

	return result
}

func (u *Sessioninfo) SetToOneReferenceID(name, ID string) error {
	if name == "category" {
		u.CategoryID = ID
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}

// SetToManyReferenceIDs sets the leafs reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (u *Sessioninfo) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "images" {
		u.ImagesIDs = IDs
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs adds some new leafs that a users loves so much
func (u *Sessioninfo) AddToManyIDs(name string, IDs []string) error {
	if name == "images" {
		u.ImagesIDs = append(u.ImagesIDs, IDs...)
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// DeleteToManyIDs removes some leafs from a users because they made him very sick
func (u *Sessioninfo) DeleteToManyIDs(name string, IDs []string) error {
	if name == "images" {
		for _, ID := range IDs {
			for pos, oldID := range u.ImagesIDs {
				if ID == oldID {
					// match, this ID must be removed
					u.ImagesIDs = append(u.ImagesIDs[:pos], u.ImagesIDs[pos+1:]...)
				}
			}
		}
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

func (u *Sessioninfo) GetConditionsBsonM(parameters map[string][]string) bson.M {
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
