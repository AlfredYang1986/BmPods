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

// BmStudentStorage stores all applys
type BmStudentStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmStudentStorage) NewStudentStorage(args []BmDaemons.BmDaemon) *BmStudentStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmStudentStorage{mdb}
}

// GetAll returns the model map (because we need the ID as key too)
func (s BmStudentStorage) GetAll(r api2go.Request, skip int, take int) []*BmModel.Student {
	in := BmModel.Student{}
	var out []BmModel.Student
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.Student
		for i := 0; i < len(out); i++ {
			ptr := out[i]
			s.db.ResetIdWithId_(&ptr)
			tmp = append(tmp, &ptr)
		}
		return tmp
	} else {
		return nil //make(map[string]*BmModel.Student)
	}
}

// GetOne model
func (s BmStudentStorage) GetOne(id string) (BmModel.Student, error) {
	in := BmModel.Student{ID: id}
	out := BmModel.Student{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Student for id %s not found", id)
	return BmModel.Student{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a model
func (s *BmStudentStorage) Insert(c BmModel.Student) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmStudentStorage) Delete(id string) error {
	in := BmModel.Student{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Student with id %s does not exist", id)
	}

	return nil
}

// Update a model
func (s *BmStudentStorage) Update(c BmModel.Student) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Student with id does not exist")
	}

	return nil
}

func (s *BmStudentStorage) Count(req api2go.Request, c BmModel.Student) int {
	r, _ := s.db.Count(req, &c)
	return r
}
