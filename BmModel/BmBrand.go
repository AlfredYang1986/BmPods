package BmModel

import (
	"errors"

	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
)

type Brand struct {
	ID  string        `json:"id"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	Title      string   `json:"title" bson:"title"`
	SubTitle   string   `json:"sub-title" bson:"sub-title"`
	Found      float64  `json:"found" bson:"found"`
	FoundStory string   `json:"found-story" bson:"found-story"`
	Logo       string   `json:"logo" bson:"logo"`             //品牌logo
	Slogan     string   `json:"slogan" bson:"slogan"`         //一句话介绍
	BrandTags  []string `json:"brand-tags" bson:"brand-tags"` //HightLight[与众不同],3-5条,一条5个字
	EduIdea    string   `json:"edu-idea" bson:"edu-idea"`     //教育理念
	AboutUs    string   `json:"about-us" bson:"about-us"`     //团队

	CategoryID string    `json:"-" bson:"category-id"`
	Cat        *Category `json:"-"`
	ImagesIDs  []string  `json:"-" bson:"image-ids"`
	Imgs       []*Image  `json:"-"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Brand) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Brand) SetID(id string) error {
	u.ID = id
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Brand) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "images",
			Name: "images",
		},
		{
			Type: "categories",
			Name: "category",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Brand) GetReferencedIDs() []jsonapi.ReferenceID {
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
func (u Brand) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	for key := range u.Imgs {
		result = append(result, u.Imgs[key])
	}

	if u.CategoryID != "" {
		result = append(result, u.Cat)
	}

	return result
}

func (u *Brand) SetToOneReferenceID(name, ID string) error {
	if name == "category" {
		u.CategoryID = ID
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}

// SetToManyReferenceIDs sets the leafs reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (u *Brand) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "images" {
		u.ImagesIDs = IDs
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs adds some new leafs that a users loves so much
func (u *Brand) AddToManyIDs(name string, IDs []string) error {
	if name == "images" {
		u.ImagesIDs = append(u.ImagesIDs, IDs...)
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// DeleteToManyIDs removes some leafs from a users because they made him very sick
func (u *Brand) DeleteToManyIDs(name string, IDs []string) error {
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

func (u *Brand) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
