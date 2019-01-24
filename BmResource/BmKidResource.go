package BmResource

import (
	"errors"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
)

type BmKidResource struct {
	BmKidStorage   *BmDataStorage.BmKidStorage
	BmApplyStorage *BmDataStorage.BmApplyStorage
}

func (c BmKidResource) NewKidResource(args []BmDataStorage.BmStorage) BmKidResource {
	var us *BmDataStorage.BmApplyStorage
	var cs *BmDataStorage.BmKidStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmApplyStorage" {
			us = arg.(*BmDataStorage.BmApplyStorage)
		} else if tp.Name() == "BmKidStorage" {
			cs = arg.(*BmDataStorage.BmKidStorage)
		}
	}
	return BmKidResource{BmApplyStorage: us, BmKidStorage: cs}
}

// FindAll kids
func (c BmKidResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	kidsID, ok := r.QueryParams["kidsID"]
	kids := c.BmKidStorage.GetAll(r)
	if ok {
		modelID := kidsID[0]
		filteredLeafs := []BmModel.Kid{}
		model, err := c.BmApplyStorage.GetOne(modelID)
		if err != nil {
			return &Response{}, err
		}
		for _, modelLeafID := range model.KidsIDs {
			sweet, err := c.BmKidStorage.GetOne(modelLeafID)
			if err != nil {
				return &Response{}, err
			}
			filteredLeafs = append(filteredLeafs, sweet)
		}

		return &Response{Res: filteredLeafs}, nil
	}
	return &Response{Res: kids}, nil
}

// FindOne choc
func (c BmKidResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.BmKidStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new choc
func (c BmKidResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Kid)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := c.BmKidStorage.Insert(choc)
	choc.ID = id
	return &Response{Res: choc, Code: http.StatusCreated}, nil
}

// Delete a choc :(
func (c BmKidResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.BmKidStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a choc
func (c BmKidResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Kid)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := c.BmKidStorage.Update(choc)
	return &Response{Res: choc, Code: http.StatusNoContent}, err
}
