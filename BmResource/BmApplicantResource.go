package BmResource

import (
	"errors"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
)

type BmApplicantResource struct {
	BmApplicantStorage *BmDataStorage.BmApplicantStorage
	BmApplyStorage     *BmDataStorage.BmApplyStorage
	BmKidStorage       *BmDataStorage.BmKidStorage
}

func (c BmApplicantResource) NewApplicantResource(args []BmDataStorage.BmStorage) BmApplicantResource {
	var ks *BmDataStorage.BmKidStorage
	var as *BmDataStorage.BmApplicantStorage
	var apys *BmDataStorage.BmApplyStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmApplicantStorage" {
			as = arg.(*BmDataStorage.BmApplicantStorage)
		} else if tp.Name() == "BmApplyStorage" {
			apys = arg.(*BmDataStorage.BmApplyStorage)
		} else if tp.Name() == "BmKidStorage" {
			ks = arg.(*BmDataStorage.BmKidStorage)
		}
	}
	return BmApplicantResource{BmApplicantStorage: as, BmApplyStorage: apys, BmKidStorage: ks}
}

func (c BmApplicantResource) FindAll(r api2go.Request) (api2go.Responder, error) {

	appliesID, ok := r.QueryParams["appliesID"]
	if ok {
		modelRootID := appliesID[0]
		modelRoot, err := c.BmApplyStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.ApplicantID
		if modelID != "" {
			model, err := c.BmApplicantStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	kidsID, ok := r.QueryParams["kidsID"]
	if ok {
		modelRootID := kidsID[0]
		modelRoot, err := c.BmKidStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.ApplicantID
		if modelID != "" {
			model, err := c.BmApplicantStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	result := c.BmApplicantStorage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

func (s BmApplicantResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	result := []BmModel.Applicant{}
	return 100, &Response{Res: result}, nil
}

// FindOne ape
func (c BmApplicantResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.BmApplicantStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new ape
func (c BmApplicantResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	ape, ok := obj.(BmModel.Applicant)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := c.BmApplicantStorage.Insert(ape)
	ape.ID = id
	return &Response{Res: ape, Code: http.StatusCreated}, nil
}

// Delete a ape :(
func (c BmApplicantResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.BmApplicantStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a ape
func (c BmApplicantResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	ape, ok := obj.(BmModel.Applicant)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := c.BmApplicantStorage.Update(ape)
	return &Response{Res: ape, Code: http.StatusNoContent}, err
}
