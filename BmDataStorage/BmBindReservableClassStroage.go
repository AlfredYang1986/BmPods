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

// BmBindReservableClassStorage stores all of the tasty modelleaf, needs to be injected into
// BindReservableClass and BindReservableClass Resource. In the real world, you would use a database for that.
type BmBindReservableClassStorage struct {
	bindReservableClasses  map[string]*BmModel.BindReservableClass
	idCount int

	db *BmMongodb.BmMongodb
}

func (s BmBindReservableClassStorage) NewBindReservableClassStorage(args []BmDaemons.BmDaemon) *BmBindReservableClassStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmBindReservableClassStorage{make(map[string]*BmModel.BindReservableClass), 1, mdb}
}

// GetAll of the modelleaf
func (s BmBindReservableClassStorage) GetAll(r api2go.Request) []BmModel.BindReservableClass {
	in := BmModel.BindReservableClass{}
	out := []BmModel.BindReservableClass{}
	err := s.db.FindMulti(r, &in, &out, -1, -1)
	if err == nil {
		for i, iter := range out {
			s.db.ResetIdWithId_(&iter)
			out[i] = iter
		}
		return out
	} else {
		return nil
	}
}

// GetOne tasty modelleaf
func (s BmBindReservableClassStorage) GetOne(id string) (BmModel.BindReservableClass, error) {
	in := BmModel.BindReservableClass{ID: id}
	out := BmModel.BindReservableClass{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("BindReservableClass for id %s not found", id)
	return BmModel.BindReservableClass{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *BmBindReservableClassStorage) Insert(c BmModel.BindReservableClass) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmBindReservableClassStorage) Delete(id string) error {
	in := BmModel.BindReservableClass{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("BindReservableClass with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing modelleaf
func (s *BmBindReservableClassStorage) Update(c BmModel.BindReservableClass) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("BindReservableClass with id does not exist")
	}

	return nil
}
