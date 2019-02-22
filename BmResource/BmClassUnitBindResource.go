package BmResource

import (
	"errors"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
)

type BmClassUnitBindResource struct {
	BmClassUnitBindStorage     *BmDataStorage.BmClassUnitBindStorage
	BmClassStorage    		   *BmDataStorage.BmClassStorage
	BmUnitStorage        	   *BmDataStorage.BmUnitStorage
}

func (s BmClassUnitBindResource) NewClassUnitBindResource(args []BmDataStorage.BmStorage) BmClassUnitBindResource {
	var ds *BmDataStorage.BmClassUnitBindStorage
	var ts *BmDataStorage.BmClassStorage
	var us *BmDataStorage.BmUnitStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmClassUnitBindStorage" {
			ds = arg.(*BmDataStorage.BmClassUnitBindStorage)
		} else if tp.Name() == "BmClassStorage" {
			ts = arg.(*BmDataStorage.BmClassStorage)
		}else if tp.Name() == "BmUnitStorage" {
			us = arg.(*BmDataStorage.BmUnitStorage)
		}
	}
	return BmClassUnitBindResource{BmClassUnitBindStorage: ds, BmClassStorage: ts,BmUnitStorage: us}
}

// FindAll to satisfy api2go data source interface
func (s BmClassUnitBindResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []*BmModel.ClassUnitBind
	result = s.BmClassUnitBindStorage.GetAll(r, -1, -1)

	return &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the user with the given ID, otherwise an error
func (s BmClassUnitBindResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmClassUnitBindStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	return &Response{Res: model}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BmClassUnitBindResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.ClassUnitBind)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := s.BmClassUnitBindStorage.Insert(model)
	model.ID = id
	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s BmClassUnitBindResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.BmClassUnitBindStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the user
func (s BmClassUnitBindResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	user, ok := obj.(BmModel.ClassUnitBind)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BmClassUnitBindStorage.Update(user)
	return &Response{Res: user, Code: http.StatusNoContent}, err
}
