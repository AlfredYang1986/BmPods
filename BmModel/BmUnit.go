package BmModel

import (
	"errors"
	"github.com/manyminds/api2go/jsonapi"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"strconv"
)

type Unit struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	Status     float64 `json:"status" bson:"status"`
	StartDate  float64 `json:"start-date" bson:"start-date"`
	EndDate    float64 `json:"end-date" bson:"end-date"`
	CourseTime float64 `json:"course-time" bson:"course-time"` //课时
	BrandId    string  `json:"brand-id" bson:"brand-id"`

	TeacherID string  `json:"teacher-id" bson:"teacher-id"`
	Teacher   Teacher `json:"-"`
	ClassUnitBindID   string  `json:"classunitbind-id" bson:"classunitbind-id"`
	ClassUnitBind     ClassUnitBind   `json:"-"`
	//通过room过滤unit
	RoomID string `json:"room-id" bson:"room-id"`
	Room   Room   `json:"-"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (u Unit) GetID() string {
	return u.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (u *Unit) SetID(id string) error {
	u.ID = id
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (u Unit) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{
		{
			Type: "teachers",
			Name: "teacher",
		},
		{
			Type: "classunitbinds",
			Name: "classunitbind",
		},
		{
			Type: "rooms",
			Name: "room",
		},
	}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (u Unit) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}

	if u.TeacherID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.TeacherID,
			Type: "teachers",
			Name: "teacher",
		})
	}
	if u.ClassUnitBindID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.ClassUnitBindID,
			Type: "classunitbinds",
			Name: "classunitbind",
		})
	}

	if u.RoomID != "" {
		result = append(result, jsonapi.ReferenceID{
			ID:   u.RoomID,
			Type: "rooms",
			Name: "room",
		})
	}

	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (u Unit) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}

	if u.TeacherID != "" {
		result = append(result, u.Teacher)
	}
	if u.ClassUnitBindID != "" {
		result = append(result, u.ClassUnitBind)
	}

	if u.RoomID != "" {
		result = append(result, u.Room)
	}

	// TODO: sessionable

	return result
}

func (u *Unit) SetToOneReferenceID(name, ID string) error {
	if name == "teacher" {
		u.TeacherID = ID
		return nil
	}
	if name == "classunitbind" {
		u.ClassUnitBindID = ID
		return nil
	}

	if name == "room" {
		u.RoomID = ID
		return nil
	}

	return errors.New("There is no to-one relationship with the name " + name)
}

func (u *Unit) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	for k, v := range parameters {
		switch k {
		case "brand-id":
			rst[k] = v[0]
		case "room-id":
			rst[k] = v[0]
		case "status":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			rst[k] = val
		case "lt[start-date]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$lt"] = val
			rst["start-date"] = r
		case "lte[start-date]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$lte"] = val
			rst["start-date"] = r
		case "gt[start-date]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$gt"] = val
			rst["start-date"] = r
		case "gte[start-date]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$gte"] = val
			rst["start-date"] = r
		case "lt[end-date]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$lt"] = val
			rst["end-date"] = r
		case "lte[end-date]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$lte"] = val
			rst["end-date"] = r
		case "gt[end-date]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$gt"] = val
			rst["end-date"] = r
		case "gte[end-date]":
			val, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic(err.Error())
			}
			r := make(map[string]interface{})
			r["$gte"] = val
			rst["end-date"] = r
		}
	}

	return rst
}

type Units []*Unit

func (bd Units) SortByStartDate(increasing bool) error {
	courseUnits := bd
	if courseUnits == nil {
		return nil
	}
	if increasing {
		sort.Sort(BmUnitsWrapper{courseUnits, func(cu1, cu2 *Unit) bool {
			return cu1.StartDate < cu2.StartDate //按开始时间递增排序
		}})
	} else {
		sort.Sort(BmUnitsWrapper{courseUnits, func(cu1, cu2 *Unit) bool {
			return cu1.StartDate > cu2.StartDate //按开始时间递减排序
		}})
	}

	return nil
}

func (bd Units) SortByEndDate(increasing bool) error {
	courseUnits := bd
	if courseUnits == nil {
		return nil
	}
	if increasing {
		sort.Sort(BmUnitsWrapper{courseUnits, func(cu1, cu2 *Unit) bool {
			return cu1.EndDate < cu2.EndDate //按结束时间递增排序
		}})
	} else {
		sort.Sort(BmUnitsWrapper{courseUnits, func(cu1, cu2 *Unit) bool {
			return cu1.EndDate > cu2.EndDate //按结束时间递减排序
		}})
	}

	return nil
}

type BmUnitsWrapper struct {
	units  []*Unit
	sortBy func(cu1, cu2 *Unit) bool
}

func (bd BmUnitsWrapper) Len() int {
	return len(bd.units)
}

func (bd BmUnitsWrapper) Swap(i, j int) {
	bd.units[i], bd.units[j] = bd.units[j], bd.units[i]
}

func (bd BmUnitsWrapper) Less(i, j int) bool {
	return bd.sortBy(bd.units[i], bd.units[j])
}
