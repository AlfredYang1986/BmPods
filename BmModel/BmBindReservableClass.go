package BmModel

import "gopkg.in/mgo.v2/bson"

// BindReservableClass is the BindReservableClass that a user consumes in order to get fat and happy
type BindReservableClass struct {
	ID               string        `json:"-"`
	Id_              bson.ObjectId `json:"-" bson:"_id"`
	ReservableitemId string        `json:"reservableitem-id" bson:"reservableitem-id"`
	ClassId          string        `json:"class-id" bson:"class-id"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (c BindReservableClass) GetID() string {
	return c.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (c *BindReservableClass) SetID(id string) error {
	c.ID = id
	return nil
}

func (u *BindReservableClass) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}

func (u * BindReservableClass) GetConditionsWithClassId() bson.M {
	return bson.M{"class-id": u.ClassId}
}

func (u * BindReservableClass) GetConditionsWithResId() bson.M {
	return bson.M{"reservableitem-id": u.ReservableitemId}
}

func (u * BindReservableClass) GetConditionsWithResIdAndClsId() bson.M {
	return bson.M{
		"reservableitem-id": u.ReservableitemId,
		"class-id": u.ClassId}
}