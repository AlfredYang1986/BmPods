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

// CategoryStorage stores all users
type BmCategoryStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmCategoryStorage) NewCategoryStorage(args []BmDaemons.BmDaemon) *BmCategoryStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmCategoryStorage{mdb}
}

// GetAll returns the user map (because we need the ID as key too)
func (s BmCategoryStorage) GetAll(r api2go.Request, skip int, take int) []*BmModel.Category {
	in := BmModel.Category{}
	var out []BmModel.Category
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.Category
		//tmp := make(map[string]*BmModel.Category)
		for _, iter := range out {
			s.db.ResetIdWithId_(&iter)
			tmp = append(tmp, &iter)
			//tmp[iter.ID] = &iter
		}
		return tmp
	} else {
		return nil //make(map[string]*BmModel.Category)
	}
}

// GetOne user
func (s BmCategoryStorage) GetOne(id string) (BmModel.Category, error) {
	in := BmModel.Category{ID: id}
	out := BmModel.Category{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Category for id %s not found", id)
	return BmModel.Category{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a user
func (s *BmCategoryStorage) Insert(c BmModel.Category) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmCategoryStorage) Delete(id string) error {
	in := BmModel.Category{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Category with id %s does not exist", id)
	}

	return nil
}

// Update a user
func (s *BmCategoryStorage) Update(c BmModel.Category) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Category with id does not exist")
	}

	return nil
}

func (s *BmCategoryStorage) Count(req api2go.Request, c BmModel.Category) int {
	r, _ := s.db.Count(req, &c)
	return r
}
