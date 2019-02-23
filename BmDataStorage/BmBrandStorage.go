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

// BmBrandStorage stores all applys
type BmBrandStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmBrandStorage) NewBrandStorage(args []BmDaemons.BmDaemon) *BmBrandStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmBrandStorage{mdb}
}

// GetAll returns the model map (because we need the ID as key too)
func (s BmBrandStorage) GetAll(r api2go.Request, skip int, take int) []*BmModel.Brand {
	in := BmModel.Brand{}
	var out []BmModel.Brand
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.Brand
		for i := 0; i < len(out); i++ {
			ptr := out[i]
			s.db.ResetIdWithId_(&ptr)
			tmp = append(tmp, &ptr)
		}
		return tmp
	} else {
		return nil //make(map[string]*BmModel.Brand)
	}
}

// GetOne model
func (s BmBrandStorage) GetOne(id string) (BmModel.Brand, error) {
	in := BmModel.Brand{ID: id}
	out := BmModel.Brand{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Brand for id %s not found", id)
	return BmModel.Brand{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a model
func (s *BmBrandStorage) Insert(c BmModel.Brand) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmBrandStorage) Delete(id string) error {
	in := BmModel.Brand{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Brand with id %s does not exist", id)
	}

	return nil
}

// Update a model
func (s *BmBrandStorage) Update(c BmModel.Brand) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Brand with id does not exist")
	}

	return nil
}

func (s *BmBrandStorage) Count(req api2go.Request, c BmModel.Brand) int {
	r, _ := s.db.Count(req, &c)
	return r
}
