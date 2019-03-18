package BmDataStorage

import (
	"errors"
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
)

// BmAttachableStorage stores all of the tasty modelleaf, needs to be injected into
// Attachable Resource. In the real world, you would use a database for that.
type BmAttachableStorage struct {
	attachables    map[string]*BmModel.Attachable
	idCount int

	db *BmMongodb.BmMongodb
}

func (s BmAttachableStorage) NewAttachableStorage(args []BmDaemons.BmDaemon) *BmAttachableStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmAttachableStorage{make(map[string]*BmModel.Attachable), 1, mdb}
}

// GetAll of the modelleaf
func (s BmAttachableStorage) GetAll(r api2go.Request, skip int, take int) []BmModel.Attachable {
	in := BmModel.Attachable{}
	out := []BmModel.Attachable{}
	err := s.db.FindMulti(r, &in, &out, skip, take)
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
func (s BmAttachableStorage) GetOne(id string) (BmModel.Attachable, error) {
	in := BmModel.Attachable{ID: id}
	out := BmModel.Attachable{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Attachable for id %s not found", id)
	return BmModel.Attachable{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *BmAttachableStorage) Insert(c BmModel.Attachable) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmAttachableStorage) Delete(id string) error {
	in := BmModel.Attachable{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Attachable with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing modelleaf
func (s *BmAttachableStorage) Update(c BmModel.Attachable) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Attachable with id does not exist")
	}

	return nil
}
