package BmDataStorage

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmPods/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
)

//BmYardStorage stores all applys
type BmYardStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmYardStorage) NewYardStorage(args []BmDaemons.BmDaemon) *BmYardStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmYardStorage{mdb}
}

func (s BmYardStorage) GetAll(r api2go.Request, skip, take int) []*BmModel.Yard {
	in := BmModel.Yard{}

	var out []BmModel.Yard
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.Yard
		for i := 0; i < len(out); i++ {
			ptr := out[i]
			s.db.ResetIdWithId_(&ptr)
			tmp = append(tmp, &ptr)
		}
		return tmp
	} else {
		return nil
	}
}

// GetOne model
func (s BmYardStorage) GetOne(id string) (BmModel.Yard, error) {
	in := BmModel.Yard{ID: id}
	out := BmModel.Yard{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Yard for id %s not found", id)
	return BmModel.Yard{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a model
func (s *BmYardStorage) Insert(c BmModel.Yard) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmYardStorage) Delete(id string) error {
	in := BmModel.Yard{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Yard with id %s does not exist", id)
	}

	return nil
}

// Update a model
func (s *BmYardStorage) Update(c BmModel.Yard) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Yard with id does not exist")
	}

	return nil
}

// Count a model
func (s *BmYardStorage) Count(req api2go.Request, c BmModel.Yard) int {
	r, _ := s.db.Count(req, &c)
	return r
}
