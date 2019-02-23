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

// BmTeacherStorage stores all of the tasty chocolate, needs to be injected into
// Teacher and Teacher Resource. In the real world, you would use a database for that.
type BmTeacherStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmTeacherStorage) NewTeacherStorage(args []BmDaemons.BmDaemon) *BmTeacherStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmTeacherStorage{mdb}
}

// GetAll of the chocolate
func (s BmTeacherStorage) GetAll(r api2go.Request, skip int, take int) []BmModel.Teacher {
	in := BmModel.Teacher{}
	out := []BmModel.Teacher{}
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

// GetOne tasty chocolate
func (s BmTeacherStorage) GetOne(id string) (BmModel.Teacher, error) {
	in := BmModel.Teacher{ID: id}
	out := BmModel.Teacher{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Teacher for id %s not found", id)
	return BmModel.Teacher{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *BmTeacherStorage) Insert(c BmModel.Teacher) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmTeacherStorage) Delete(id string) error {
	in := BmModel.Teacher{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Teacher with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing chocolate
func (s *BmTeacherStorage) Update(c BmModel.Teacher) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Teacher with id does not exist")
	}

	return nil
}

func (s *BmTeacherStorage) Count(req api2go.Request, c BmModel.Teacher) int {
	r, _ := s.db.Count(req, &c)
	return r
}
