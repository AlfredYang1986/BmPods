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

// BmKidStorage stores all of the tasty modelleaf, needs to be injected into
// Kid Resource. In the real world, you would use a database for that.
type BmKidStorage struct {
	kids    map[string]*BmModel.Kid
	idCount int

	db *BmMongodb.BmMongodb
}

func (s BmKidStorage) NewKidStorage(args []BmDaemons.BmDaemon) *BmKidStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmKidStorage{make(map[string]*BmModel.Kid), 1, mdb}
}

// GetAll of the modelleaf
func (s BmKidStorage) GetAll(r api2go.Request) []BmModel.Kid {
	in := BmModel.Kid{}
	out := []BmModel.Kid{}
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
func (s BmKidStorage) GetOne(id string) (BmModel.Kid, error) {
	in := BmModel.Kid{ID: id}
	out := BmModel.Kid{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Kid for id %s not found", id)
	return BmModel.Kid{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *BmKidStorage) Insert(c BmModel.Kid) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmKidStorage) Delete(id string) error {
	in := BmModel.Kid{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Kid with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing modelleaf
func (s *BmKidStorage) Update(c BmModel.Kid) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Kid with id does not exist")
	}

	return nil
}
