package BmDataStorage

import (
	"errors"
	"fmt"
	"github.com/alfredyang1986/BmPods/BmDaemons"
	"github.com/alfredyang1986/BmPods/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
)

// BmUnitStorage stores all of the tasty chocolate, needs to be injected into
// User and Unit Resource. In the real world, you would use a database for that.
type BmUnitStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmUnitStorage) NewUnitStorage(args []BmDaemons.BmDaemon) *BmUnitStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmUnitStorage{mdb}
}

// GetAll of the chocolate
func (s BmUnitStorage) GetAll(skip int, take int) []*BmModel.Unit {
	in := BmModel.Unit{}
	var out []BmModel.Unit
	err := s.db.FindMulti(&in, &out, skip, take)
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
		return out, nil
	}
	errMessage := fmt.Sprintf("User for id %s not found", id)
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

func (s *BmUnitStorage) Count(c BmModel.Unit) int {
	r, _ := s.db.Count(&c)
	return r
}
