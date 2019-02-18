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

// CategoryStorage stores all users
type BmNewCategoryStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmNewCategoryStorage) NewNewCategoryStorage(args []BmDaemons.BmDaemon) *BmNewCategoryStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmNewCategoryStorage{mdb}
}

// GetAll returns the user map (because we need the ID as key too)
func (s BmNewCategoryStorage) GetAll(r api2go.Request, skip int, take int) []*BmModel.NewCategory {
	in := BmModel.NewCategory{}
	//out := []BmModel.Kid{}
	var out []BmModel.NewCategory
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.NewCategory
		for _, iter := range out {
			s.db.ResetIdWithId_(&iter)
			tmpIter:=iter
			tmp = append(tmp, &tmpIter)
		}
		return tmp
	} else {
		return nil //make(map[string]*BmModel.Category)
	}
}

// GetOne user
func (s BmNewCategoryStorage) GetOne(id string) (BmModel.NewCategory, error) {
	in := BmModel.NewCategory{ID: id}
	out := BmModel.NewCategory{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Category for id %s not found", id)
	return BmModel.NewCategory{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a user
func (s *BmNewCategoryStorage) Insert(c BmModel.NewCategory) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}
	return tmp
}

// Delete one :(
func (s *BmNewCategoryStorage) Delete(id string) error {
	in := BmModel.NewCategory{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Category with id %s does not exist", id)
	}

	return nil
}

// Update a user
func (s *BmNewCategoryStorage) Update(c BmModel.NewCategory) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("NewCategory with id does not exist")
	}

	return nil
}

func (s *BmNewCategoryStorage) Count(req api2go.Request, c BmModel.NewCategory) int {
	r, _ := s.db.Count(req, &c)
	return r
}
