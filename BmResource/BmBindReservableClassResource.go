package BmResource

import (
	"errors"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
)

type BmBindReservableClassResource struct {
	BmBindReservableClassStorage       *BmDataStorage.BmBindReservableClassStorage
}

func (c BmBindReservableClassResource) NewBindReservableClassResource(args []BmDataStorage.BmStorage) BmBindReservableClassResource {
	var brc *BmDataStorage.BmBindReservableClassStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmBindReservableClassStorage" {
			brc = arg.(*BmDataStorage.BmBindReservableClassStorage)
		}
	}
	return BmBindReservableClassResource{BmBindReservableClassStorage: brc}
}

func (c BmBindReservableClassResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	result := []BmModel.BindReservableClass{}
	result = c.BmBindReservableClassStorage.GetAll(r)
	return &Response{Res: result}, nil
}

// FindOne choc
func (c BmBindReservableClassResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.BmBindReservableClassStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new choc
func (c BmBindReservableClassResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.BindReservableClass)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := c.BmBindReservableClassStorage.Insert(choc)
	choc.ID = id
	return &Response{Res: choc, Code: http.StatusCreated}, nil
}

// Delete a choc :(
func (c BmBindReservableClassResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.BmBindReservableClassStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a choc
func (c BmBindReservableClassResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.BindReservableClass)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := c.BmBindReservableClassStorage.Update(choc)
	return &Response{Res: choc, Code: http.StatusNoContent}, err
}
