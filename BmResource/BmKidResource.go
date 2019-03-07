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
	BmStudentStorage *BmDataStorage.BmStudentStorage
}

func (c BmKidResource) NewKidResource(args []BmDataStorage.BmStorage) BmKidResource {
	var us *BmDataStorage.BmApplyStorage
	var cs *BmDataStorage.BmKidStorage
	var ss *BmDataStorage.BmStudentStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmApplyStorage" {
			us = arg.(*BmDataStorage.BmApplyStorage)
		} else if tp.Name() == "BmKidStorage" {
			cs = arg.(*BmDataStorage.BmKidStorage)
		} else if tp.Name() == "BmStudentStorage" {
			ss = arg.(*BmDataStorage.BmStudentStorage)
		}
	}
	return BmKidResource{BmApplyStorage: us, BmKidStorage: cs, BmStudentStorage: ss}
}

// FindAll kids
func (c BmKidResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	appliesID, ok := r.QueryParams["appliesID"]
	if ok {
		modelRootID := appliesID[0]
		modelRoot, err := c.BmApplyStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		if len(modelRoot.KidsIDs) != 0 {
			r.QueryParams["kidsids"] = modelRoot.KidsIDs
			models := c.BmKidStorage.GetAll(r)
			return &Response{Res: models}, nil
		}
		return &Response{}, nil
	}

	studentsID, ok := r.QueryParams["studentsID"]
	if ok {
		modelRootID := studentsID[0]
		modelRoot, err := c.BmStudentStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.KidID
		if modelID != "" {
			model, err := c.BmKidStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
		return &Response{}, nil
	}

	kids := c.BmKidStorage.GetAll(r)
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
	panic("不可删除")
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
