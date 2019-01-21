package BmResource

import (
	"errors"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
)

type BmImageResource struct {
	BmImageStorage       *BmDataStorage.BmImageStorage
	BmSessioninfoStorage *BmDataStorage.BmSessioninfoStorage
}

func (c BmImageResource) NewImageResource(args []BmDataStorage.BmStorage) BmImageResource {
	var us *BmDataStorage.BmSessioninfoStorage
	var cs *BmDataStorage.BmImageStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmSessioninfoStorage" {
			us = arg.(*BmDataStorage.BmSessioninfoStorage)
		} else if tp.Name() == "BmImageStorage" {
			cs = arg.(*BmDataStorage.BmImageStorage)
		}
	}
	return BmImageResource{BmSessioninfoStorage: us, BmImageStorage: cs}
}

// FindAll images
func (c BmImageResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	sessioninfosID, ok := r.QueryParams["sessioninfosID"]
	result := []BmModel.Image{}
	if ok {
		// this means that we want to show all images of a modelRoot, this is the route
		// /v0/models/1/images
		modelRootID := sessioninfosID[0]
		// filter out images with modelRootID, in real world, you would just run a different database query

		modelRoot, err := c.BmSessioninfoStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		for _, modelID := range modelRoot.ImagesIDs {
			model, err := c.BmImageStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			result = append(result, model)
		}

		return &Response{Res: result}, nil
	}
	result = c.BmImageStorage.GetAll(r)
	return &Response{Res: result}, nil
}

// FindOne choc
func (c BmImageResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.BmImageStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new choc
func (c BmImageResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Image)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := c.BmImageStorage.Insert(choc)
	choc.ID = id
	return &Response{Res: choc, Code: http.StatusCreated}, nil
}

// Delete a choc :(
func (c BmImageResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.BmImageStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a choc
func (c BmImageResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Image)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := c.BmImageStorage.Update(choc)
	return &Response{Res: choc, Code: http.StatusNoContent}, err
}
