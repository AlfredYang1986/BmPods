package BmResource

import (
	"errors"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
)

type BmCatenodeResource struct {
	CatenodeStorage *BmDataStorage.BmCatenodeStorage
	BmBrandStorage     *BmDataStorage.BmBrandStorage
}

func (c BmCatenodeResource) NewCatenodeResource(args []BmDataStorage.BmStorage) BmCatenodeResource {
	var as *BmDataStorage.BmCatenodeStorage
	var bs *BmDataStorage.BmBrandStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmCatenodeStorage" {
			as = arg.(*BmDataStorage.BmCatenodeStorage)
		} else if tp.Name() == "BmBrandStorage" {
			bs = arg.(*BmDataStorage.BmBrandStorage)
		}
	}
	return BmCatenodeResource{CatenodeStorage: as, BmBrandStorage: bs}
}

// FindAll apeolates
func (c BmCatenodeResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	brandsID, ok := r.QueryParams["brandsID"]
	if ok {
		modelRootID := brandsID[0]

		modelRoot, err := c.BmBrandStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.CategoryID
		if modelID != "" {
			model, err := c.CatenodeStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			//result = append(result, model)

			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}
	result := c.CatenodeStorage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

func (c BmCatenodeResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	result := []BmModel.Catenode{}
	return 100, &Response{Res: result}, nil
}

// FindOne ape
func (c BmCatenodeResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.CatenodeStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new ape
func (c BmCatenodeResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	ape, ok := obj.(BmModel.Catenode)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := c.CatenodeStorage.Insert(ape)
	ape.ID = id
	return &Response{Res: ape, Code: http.StatusCreated}, nil
}

// Delete a ape :(
func (c BmCatenodeResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.CatenodeStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a ape
func (c BmCatenodeResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	ape, ok := obj.(BmModel.Catenode)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := c.CatenodeStorage.Update(ape)
	return &Response{Res: ape, Code: http.StatusNoContent}, err
}
