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

// ApplicantStorage stores all users
type BmApplicantStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmApplicantStorage) NewApplicantStorage(args []BmDaemons.BmDaemon) *BmApplicantStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmApplicantStorage{mdb}
}

// GetAll returns the user map (because we need the ID as key too)
func (s BmApplicantStorage) GetAll(r api2go.Request, skip int, take int) []*BmModel.Applicant {
	in := BmModel.Applicant{}
	var out []BmModel.Applicant
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.Applicant
		//tmp := make(map[string]*BmModel.Applicant)
		for _, iter := range out {
			s.db.ResetIdWithId_(&iter)
			tmp = append(tmp, &iter)
			//tmp[iter.ID] = &iter
		}
		return tmp
	} else {
		return nil //make(map[string]*BmModel.Applicant)
	}
}

// GetOne user
func (s BmApplicantStorage) GetOne(id string) (BmModel.Applicant, error) {
	in := BmModel.Applicant{ID: id}
	out := BmModel.Applicant{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Applicant for id %s not found", id)
	return BmModel.Applicant{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a user
func (s *BmApplicantStorage) Insert(c BmModel.Applicant) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmApplicantStorage) Delete(id string) error {
	in := BmModel.Applicant{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Applicant with id %s does not exist", id)
	}

	return nil
}

// Update a user
func (s *BmApplicantStorage) Update(c BmModel.Applicant) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Applicant with id does not exist")
	}

	return nil
}

func (s *BmApplicantStorage) Count(req api2go.Request, c BmModel.Applicant) int {
	r, _ := s.db.Count(req, &c)
	return r
}
