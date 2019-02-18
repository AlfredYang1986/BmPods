package BmResource

import (
	"errors"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
)

type BmNewCategoryResource struct {
	NewCategoryStorage *BmDataStorage.BmNewCategoryStorage
	BmBrandStorage     *BmDataStorage.BmBrandStorage
}

func (c BmNewCategoryResource) NewNewCategoryResource(args []BmDataStorage.BmStorage) BmNewCategoryResource {
	var as *BmDataStorage.BmNewCategoryStorage
	var bs *BmDataStorage.BmBrandStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmNewCategoryStorage" {
			as = arg.(*BmDataStorage.BmNewCategoryStorage)
		} else if tp.Name() == "BmBrandStorage" {
			bs = arg.(*BmDataStorage.BmBrandStorage)
		}
	}
	return BmNewCategoryResource{NewCategoryStorage: as, BmBrandStorage: bs}
}

// FindAll apeolates
func (c BmNewCategoryResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	brandsID, ok := r.QueryParams["brandsID"]
	if ok {
		modelRootID := brandsID[0]

		modelRoot, err := c.BmBrandStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.CategoryID
		if modelID != "" {
			model, err := c.NewCategoryStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			//result = append(result, model)

			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}
	result := c.NewCategoryStorage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

func (c BmNewCategoryResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	result := []BmModel.NewCategory{}
	return 100, &Response{Res: result}, nil
}

// FindOne ape
func (c BmNewCategoryResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.NewCategoryStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new ape
func (c BmNewCategoryResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	ape, ok := obj.(BmModel.NewCategory)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := c.NewCategoryStorage.Insert(ape)
	ape.ID = id
	return &Response{Res: ape, Code: http.StatusCreated}, nil
}

// Delete a ape :(
func (c BmNewCategoryResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.NewCategoryStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a ape
func (c BmNewCategoryResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	ape, ok := obj.(BmModel.NewCategory)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := c.NewCategoryStorage.Update(ape)
	return &Response{Res: ape, Code: http.StatusNoContent}, err
}
