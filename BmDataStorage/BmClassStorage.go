package BmDataStorage

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alfredyang1986/BmPods/BmDaemons"
	"github.com/alfredyang1986/BmPods/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
)

// BmClassStorage stores all classes
type BmClassStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmClassStorage) NewClassStorage(args []BmDaemons.BmDaemon) *BmClassStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmClassStorage{mdb}
}

// GetAll returns the model map (because we need the ID as key too)
func (s BmClassStorage) GetAll(skip int, take int) []*BmModel.Class {
	in := BmModel.Class{}
	var out []BmModel.Class
	err := s.db.FindMulti(&in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.Class
		for i := 0; i < len(out); i++ {
			ptr := out[i]
			s.db.ResetIdWithId_(&ptr)
			tmp = append(tmp, &ptr)
		}
		return tmp
	} else {
		return nil //make(map[string]*BmModel.Class)
	}
}

// GetOne model
func (s BmClassStorage) GetOne(id string) (BmModel.Class, error) {
	in := BmModel.Class{ID: id}
	out := BmModel.Class{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Class for id %s not found", id)
	return BmModel.Class{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a model
func (s *BmClassStorage) Insert(c BmModel.Class) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmClassStorage) Delete(id string) error {
	in := BmModel.Class{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Class with id %s does not exist", id)
	}

	return nil
}

// Update a model
func (s *BmClassStorage) Update(c BmModel.Class) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Class with id does not exist")
	}

	return nil
}

func (s *BmClassStorage) Count(c BmModel.Class) int {
	r, _ := s.db.Count(&c)
	return r
}
