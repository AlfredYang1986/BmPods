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

// BmDutyStorage stores all of the tasty chocolate, needs to be injected into
// Duty Resource. In the real world, you would use a database for that.
type BmDutyStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmDutyStorage) NewDutyStorage(args []BmDaemons.BmDaemon) *BmDutyStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmDutyStorage{mdb}
}

// GetAll of the chocolate
func (s BmDutyStorage) GetAll(r api2go.Request, skip int, take int) []*BmModel.Duty {
	in := BmModel.Duty{}
	var out []BmModel.Duty
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.Duty
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

// GetOne tasty chocolate
func (s BmDutyStorage) GetOne(id string) (BmModel.Duty, error) {
	in := BmModel.Duty{ID: id}
	out := BmModel.Duty{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {

		if out.TeacherID != "" {
			item, err := BmTeacherStorage{db: s.db}.GetOne(out.TeacherID)
			if err != nil {
				return BmModel.Duty{}, err
			}
			out.Teacher = item
		}

		return out, nil
	}
	errMessage := fmt.Sprintf("Duty for id %s not found", id)
	return BmModel.Duty{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *BmDutyStorage) Insert(c BmModel.Duty) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmDutyStorage) Delete(id string) error {
	in := BmModel.Duty{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Duty with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing chocolate
func (s *BmDutyStorage) Update(c BmModel.Duty) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Duty with id does not exist")
	}

	return nil
}

func (s *BmDutyStorage) Count(req api2go.Request, c BmModel.Duty) int {
	r, _ := s.db.Count(req, &c)
	return r
}
