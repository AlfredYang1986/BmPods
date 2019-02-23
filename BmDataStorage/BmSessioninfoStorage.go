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

// BmSessioninfoStorage stores all sessioninfos
type BmSessioninfoStorage struct {
	sessioninfos map[string]*BmModel.Sessioninfo
	idCount      int

	db *BmMongodb.BmMongodb
}

func (s BmSessioninfoStorage) NewSessioninfoStorage(args []BmDaemons.BmDaemon) *BmSessioninfoStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmSessioninfoStorage{make(map[string]*BmModel.Sessioninfo), 1, mdb}
}

// GetAll returns the model map (because we need the ID as key too)
func (s BmSessioninfoStorage) GetAll(r api2go.Request, skip int, take int) []*BmModel.Sessioninfo {
	in := BmModel.Sessioninfo{}
	var out []BmModel.Sessioninfo
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.Sessioninfo
		for i := 0; i < len(out); i++ {
			ptr := out[i]
			s.db.ResetIdWithId_(&ptr)
			tmp = append(tmp, &ptr)
		}
		return tmp
	} else {
		return nil //make(map[string]*BmModel.Sessioninfo)
	}
}

// GetOne model
func (s BmSessioninfoStorage) GetOne(id string) (BmModel.Sessioninfo, error) {
	in := BmModel.Sessioninfo{ID: id}
	out := BmModel.Sessioninfo{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		if out.CategoryID != "" {
			cate, err := BmCategoryStorage{db: s.db}.GetOne(out.CategoryID)
			if err == nil {
				out.Category = cate
			}
		}
		return out, nil
	}
	errMessage := fmt.Sprintf("Sessioninfo for id %s not found", id)
	return BmModel.Sessioninfo{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a model
func (s *BmSessioninfoStorage) Insert(c BmModel.Sessioninfo) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmSessioninfoStorage) Delete(id string) error {
	in := BmModel.Sessioninfo{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Sessioninfo with id %s does not exist", id)
	}

	return nil
}

// Update a model
func (s *BmSessioninfoStorage) Update(c BmModel.Sessioninfo) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Sessioninfo with id does not exist")
	}

	return nil
}

func (s *BmSessioninfoStorage) Count(req api2go.Request, c BmModel.Sessioninfo) int {
	r, _ := s.db.Count(req, &c)
	return r
}
