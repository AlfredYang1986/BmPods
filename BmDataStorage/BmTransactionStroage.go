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

// BmTransactionStorage stores all of the tasty modelleaf, needs to be injected into
// Transaction Resource. In the real world, you would use a database for that.
type BmTransactionStorage struct {
	transactions    map[string]*BmModel.Transaction
	idCount int

	db *BmMongodb.BmMongodb
}

func (s BmTransactionStorage) NewTransactionStorage(args []BmDaemons.BmDaemon) *BmTransactionStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmTransactionStorage{make(map[string]*BmModel.Transaction), 1, mdb}
}

// GetAll of the modelleaf
func (s BmTransactionStorage) GetAll(r api2go.Request, skip int, take int) []BmModel.Transaction {
	in := BmModel.Transaction{}
	out := []BmModel.Transaction{}
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
func (s BmTransactionStorage) GetOne(id string) (BmModel.Transaction, error) {
	in := BmModel.Transaction{ID: id}
	out := BmModel.Transaction{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Transaction for id %s not found", id)
	return BmModel.Transaction{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *BmTransactionStorage) Insert(c BmModel.Transaction) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}
	return tmp
}

// Delete one :(
func (s *BmTransactionStorage) Delete(id string) error {
	in := BmModel.Transaction{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Transaction with id %s does not exist", id)
	}
	return nil
}

// Update updates an existing modelleaf
func (s *BmTransactionStorage) Update(c BmModel.Transaction) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Transaction with id does not exist")
	}
	return nil
}

func (s *BmTransactionStorage) Count(req api2go.Request, c BmModel.Transaction) int {
	r, _ := s.db.Count(req, &c)
	return r
}
