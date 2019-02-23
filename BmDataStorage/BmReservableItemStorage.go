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

// BmReservableitemStorage stores all reservableitems
type BmReservableitemStorage struct {
	reservableitems map[string]*BmModel.Reservableitem
	idCount         int

	db *BmMongodb.BmMongodb
}

func (s BmReservableitemStorage) NewReservableitemStorage(args []BmDaemons.BmDaemon) *BmReservableitemStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmReservableitemStorage{make(map[string]*BmModel.Reservableitem), 1, mdb}
}

// GetAll returns the model map (because we need the ID as key too)
func (s BmReservableitemStorage) GetAll(r api2go.Request, skip int, take int) []*BmModel.Reservableitem {
	in := BmModel.Reservableitem{}
	var out []BmModel.Reservableitem
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.Reservableitem
		for i := 0; i < len(out); i++ {
			ptr := out[i]
			s.db.ResetIdWithId_(&ptr)
			tmp = append(tmp, &ptr)
		}
		return tmp
	} else {
		return nil //make(map[string]*BmModel.Reservableitem)
	}
}

// GetOne model
func (s BmReservableitemStorage) GetOne(id string) (BmModel.Reservableitem, error) {
	in := BmModel.Reservableitem{ID: id}
	out := BmModel.Reservableitem{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Reservableitem for id %s not found", id)
	return BmModel.Reservableitem{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a model
func (s *BmReservableitemStorage) Insert(c BmModel.Reservableitem) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmReservableitemStorage) Delete(id string) error {
	in := BmModel.Reservableitem{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Reservableitem with id %s does not exist", id)
	}

	return nil
}

// Update a model
func (s *BmReservableitemStorage) Update(c BmModel.Reservableitem) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Reservableitem with id does not exist")
	}

	return nil
}

func (s *BmReservableitemStorage) Count(req api2go.Request, c BmModel.Reservableitem) int {
	r, _ := s.db.Count(req, &c)
	return r
}
