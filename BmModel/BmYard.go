package BmModel

import (
	"errors"

	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
)

type Yard struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	BrandId     string `json:"brand-id" bson:"brand-id"`
	Title       string `json:"title" bson:"title"`
	Cover       string `json:"cover" bson:"cover"`
	Description string `json:"description" bson:"description"`
	Around      string `json:"around" bson:"around"`

	//Address address.BmAddress `json:"address" bson:"relationships"`
	/**
	 * 在构建过程中，yard可能成为地址搜索的条件
	 */
	Province       string   `json:"province" bson:"province"`
	City           string   `json:"city" bson:"city"`
	District       string   `json:"district" bson:"district"`
	Address        string   `json:"address" bson:"address"`
	TrafficInfo    string   `json:"traffic-info" bson:"traffic-info"`
	Attribute      string   `json:"attribute" bson:"attribute"`
	Scenario       string   `json:"scenario" bson:"scenario"`
	OpenTime       string   `json:"open-time" bson:"open-time"`
	ServiceContact string   `json:"service-contact" bson:"service-contact"`
	Facilities     []string `json:"facilities" bson:"facilities"`
	//Friendly       []string                   `json:"friendly" bson:"friendly"`

	//RoomCount float64 `json:"room_count"`
	/**
	 * 在构建过程中，除了排课逻辑，不会通过query到Room
	 */
	ImagesIDs []string `json:"-" bson:"image-ids"`
	Images    []*Image `json:"-"`
	RoomsIDs  []string `json:"-" bson:"room-ids"`
	Rooms     []*Room  `json:"-"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Yard) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Yard) SetID(id string) error {
	u.ID = id
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Yard) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "images",
			Name: "images",
		},
		{
			Type: "rooms",
			Name: "rooms",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Yard) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	for _, kID := range u.ImagesIDs {
		result = append(result, jsonapi.ReferenceID{
			ID:   kID,
			Type: "images",
			Name: "images",
		})
	}

	for _, kID := range u.RoomsIDs {
		result = append(result, jsonapi.ReferenceID{
			ID:   kID,
			Type: "rooms",
			Name: "rooms",
		})
	}

	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (u Yard) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	for key := range u.Images {
		result = append(result, u.Images[key])
	}

	for key := range u.Rooms {
		result = append(result, u.Rooms[key])
	}

	return result
}

// SetToManyReferenceIDs sets the leafs reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (u *Yard) SetToManyReferenceIDs(name string, IDs []string) error {
	if name == "images" {
		u.ImagesIDs = IDs
		return nil
	}

	if name == "rooms" {
		u.RoomsIDs = IDs
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs adds some new leafs that a users loves so much
func (u *Yard) AddToManyIDs(name string, IDs []string) error {
	if name == "images" {
		u.ImagesIDs = append(u.ImagesIDs, IDs...)
		return nil
	}

	if name == "rooms" {
		u.RoomsIDs = append(u.RoomsIDs, IDs...)
		return nil
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

// DeleteToManyIDs removes some leafs from a users because they made him very sick
func (u *Yard) DeleteToManyIDs(name string, IDs []string) error {
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

	if name == "rooms" {
		for _, ID := range IDs {
			for pos, oldID := range u.RoomsIDs {
				if ID == oldID {
					// match, this ID must be removed
					u.RoomsIDs = append(u.RoomsIDs[:pos], u.RoomsIDs[pos+1:]...)
				}
			}
		}
	}

	return errors.New("There is no to-many relationship with the name " + name)
}

func (u *Yard) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
