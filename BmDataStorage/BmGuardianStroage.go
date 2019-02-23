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

// BmGuardianStorage stores all of the tasty modelleaf, needs to be injected into
// Guardian and Guardian Resource. In the real world, you would use a database for that.
type BmGuardianStorage struct {
	guardians map[string]*BmModel.Guardian
	idCount   int

	db *BmMongodb.BmMongodb
}

func (s BmGuardianStorage) NewGuardianStorage(args []BmDaemons.BmDaemon) *BmGuardianStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmGuardianStorage{make(map[string]*BmModel.Guardian), 1, mdb}
}

// GetAll of the modelleaf
func (s BmGuardianStorage) GetAll(r api2go.Request) []BmModel.Guardian {
	in := BmModel.Guardian{}
	out := []BmModel.Guardian{}
	err := s.db.FindMulti(r, &in, &out, -1, -1)
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
func (s BmGuardianStorage) GetOne(id string) (BmModel.Guardian, error) {
	in := BmModel.Guardian{ID: id}
	out := BmModel.Guardian{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Guardian for id %s not found", id)
	return BmModel.Guardian{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *BmGuardianStorage) Insert(c BmModel.Guardian) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmGuardianStorage) Delete(id string) error {
	in := BmModel.Guardian{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Guardian with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing modelleaf
func (s *BmGuardianStorage) Update(c BmModel.Guardian) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Guardian with id does not exist")
	}

	return nil
}
