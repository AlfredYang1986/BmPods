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

// BmImageStorage stores all of the tasty modelleaf, needs to be injected into
// Image and Image Resource. In the real world, you would use a database for that.
type BmImageStorage struct {
	images  map[string]*BmModel.Image
	idCount int

	db *BmMongodb.BmMongodb
}

func (s BmImageStorage) NewImageStorage(args []BmDaemons.BmDaemon) *BmImageStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmImageStorage{make(map[string]*BmModel.Image), 1, mdb}
}

// GetAll of the modelleaf
func (s BmImageStorage) GetAll(r api2go.Request) []BmModel.Image {
	in := BmModel.Image{}
	out := []BmModel.Image{}
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
func (s BmImageStorage) GetOne(id string) (BmModel.Image, error) {
	in := BmModel.Image{ID: id}
	out := BmModel.Image{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Image for id %s not found", id)
	return BmModel.Image{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *BmImageStorage) Insert(c BmModel.Image) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmImageStorage) Delete(id string) error {
	in := BmModel.Image{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Image with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing modelleaf
func (s *BmImageStorage) Update(c BmModel.Image) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Image with id does not exist")
	}

	return nil
}
