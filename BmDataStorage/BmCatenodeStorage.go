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

// CategoryStorage stores all users
type BmCatenodeStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmCatenodeStorage) NewCatenodeStorage(args []BmDaemons.BmDaemon) *BmCatenodeStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmCatenodeStorage{mdb}
}

// GetAll returns the user map (because we need the ID as key too)
func (s BmCatenodeStorage) GetAll(r api2go.Request, skip int, take int) []*BmModel.Catenode {
	in := BmModel.Catenode{}
	//out := []BmModel.Kid{}
	var out []BmModel.Catenode
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.Catenode
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
func (s BmCatenodeStorage) GetOne(id string) (BmModel.Catenode, error) {
	in := BmModel.Catenode{ID: id}
	out := BmModel.Catenode{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Category for id %s not found", id)
	return BmModel.Catenode{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a user
func (s *BmCatenodeStorage) Insert(c BmModel.Catenode) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}
	return tmp
}

// Delete one :(
func (s *BmCatenodeStorage) Delete(id string) error {
	in := BmModel.Catenode{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Category with id %s does not exist", id)
	}

	return nil
}

// Update a user
func (s *BmCatenodeStorage) Update(c BmModel.Catenode) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Catenode with id does not exist")
	}

	return nil
}

func (s *BmCatenodeStorage) Count(req api2go.Request, c BmModel.Catenode) int {
	r, _ := s.db.Count(req, &c)
	return r
}
