package BmDataStorage

import (
	"errors"
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmPods/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
)

// BmUnitStorage stores all of the tasty chocolate, needs to be injected into
// Unit Resource. In the real world, you would use a database for that.
type BmUnitStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmUnitStorage) NewUnitStorage(args []BmDaemons.BmDaemon) *BmUnitStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmUnitStorage{mdb}
}

// GetAll of the chocolate
func (s BmUnitStorage) GetAll(r api2go.Request, skip int, take int) []*BmModel.Unit {
	in := BmModel.Unit{}
	var out []BmModel.Unit
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.Unit
		for i := 0; i < len(out); i++ {
			ptr := out[i]
			s.db.ResetIdWithId_(&ptr)
			tmp = append(tmp, &ptr)
		}
		return tmp
	} else {
		return nil //make(map[string]*BmModel.Student)
	}
}

// GetOne tasty chocolate
func (s BmUnitStorage) GetOne(id string) (BmModel.Unit, error) {
	in := BmModel.Unit{ID: id}
	out := BmModel.Unit{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {

		//双重绑定
		if out.TeacherID != "" {
			item, err := BmTeacherStorage{db: s.db}.GetOne(out.TeacherID)
			if err != nil {
				return BmModel.Unit{}, err
			}
			out.Teacher = item
		}
		if out.ClassUnitBindID != "" {
			item, err := BmClassUnitBindStorage{db: s.db}.GetOne(out.ClassUnitBindID)
			if err != nil {
				return BmModel.Unit{}, err
			}
			out.ClassUnitBind = item
		}
		if out.RoomID != "" {
			item, err := BmRoomStorage{db: s.db}.GetOne(out.RoomID)
			if err != nil {
				return BmModel.Unit{}, err
			}
			out.Room = item
		}

		return out, nil
	}
	errMessage := fmt.Sprintf("Unit for id %s not found", id)
	return BmModel.Unit{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *BmUnitStorage) Insert(c BmModel.Unit) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmUnitStorage) Delete(id string) error {
	in := BmModel.Unit{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Unit with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing chocolate
func (s *BmUnitStorage) Update(c BmModel.Unit) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Unit with id does not exist")
	}

	return nil
}

func (s *BmUnitStorage) Count(req api2go.Request, c BmModel.Unit) int {
	r, _ := s.db.Count(req, &c)
	return r
}
