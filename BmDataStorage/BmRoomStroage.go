package BmDataStorage

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
)

type BmRoomStorage struct {
	rooms   map[string]*BmModel.Room
	idCount int

	db *BmMongodb.BmMongodb
}

func (s BmRoomStorage) NewRoomStorage(args []BmDaemons.BmDaemon) *BmRoomStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmRoomStorage{make(map[string]*BmModel.Room), 1, mdb}
}

// GetAll of the modelleaf
func (s BmRoomStorage) GetAll(r api2go.Request, skip int, take int) []BmModel.Room {
	in := BmModel.Room{}
	out := []BmModel.Room{}
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
func (s BmRoomStorage) GetOne(id string) (BmModel.Room, error) {
	in := BmModel.Room{ID: id}
	out := BmModel.Room{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Room for id %s not found", id)
	return BmModel.Room{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *BmRoomStorage) Insert(c BmModel.Room) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmRoomStorage) Delete(id string) error {
	in := BmModel.Room{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Room with id %s does not exist", id)
	}
	return nil
}

// Update updates an existing modelleaf
func (s *BmRoomStorage) Update(c BmModel.Room) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Room with id does not exist")
	}
	return nil
}

func (s *BmRoomStorage) Count(req api2go.Request, c BmModel.Room) int {
	r, _ := s.db.Count(req, &c)
	return r
}
