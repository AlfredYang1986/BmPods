package BmResource

import (
	"errors"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
	"strconv"
)

type BmSessioninfoResource struct {
	BmImageStorage          *BmDataStorage.BmImageStorage
	BmSessioninfoStorage    *BmDataStorage.BmSessioninfoStorage
	BmCategoryStorage       *BmDataStorage.BmCategoryStorage
	BmReservableitemStorage *BmDataStorage.BmReservableitemStorage
}

func (s BmSessioninfoResource) NewSessioninfoResource(args []BmDataStorage.BmStorage) BmSessioninfoResource {
	var us *BmDataStorage.BmSessioninfoStorage
	var ts *BmDataStorage.BmCategoryStorage
	var cs *BmDataStorage.BmImageStorage
	var rs *BmDataStorage.BmReservableitemStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmSessioninfoStorage" {
			us = arg.(*BmDataStorage.BmSessioninfoStorage)
		} else if tp.Name() == "BmImageStorage" {
			cs = arg.(*BmDataStorage.BmImageStorage)
		} else if tp.Name() == "BmCategoryStorage" {
			ts = arg.(*BmDataStorage.BmCategoryStorage)
		} else if tp.Name() == "BmReservableitemStorage" {
			rs = arg.(*BmDataStorage.BmReservableitemStorage)
		}
	}
	return BmSessioninfoResource{BmSessioninfoStorage: us,
		BmImageStorage:          cs,
		BmCategoryStorage:       ts,
		BmReservableitemStorage: rs}
}

// FindAll to satisfy api2go data source interface
func (s BmSessioninfoResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	resid, ok := r.QueryParams["reservableitemsID"]
	if ok {
		modelRootID := resid[0]

		modelRoot, err := s.BmReservableitemStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.SessioninfoID
		if modelID != "" {
			model, err := s.BmSessioninfoStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}

			model.Images = []*BmModel.Image{}
			for _, kID := range model.ImagesIDs {
				choc, err := s.BmImageStorage.GetOne(kID)
				if err != nil {
					return &Response{}, err
				}
				model.Images = append(model.Images, &choc)
			}

			if model.CategoryID != "" {
				cate, err := s.BmCategoryStorage.GetOne(model.CategoryID)
				if err != nil {
					return &Response{}, err
				}
				model.Category = cate
			}

			//result = append(result, model)

			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	var result []BmModel.Sessioninfo
	models := s.BmSessioninfoStorage.GetAll(r, -1, -1)

	for _, model := range models {
		// get all sweets for the model
		model.Images = []*BmModel.Image{}
		for _, kID := range model.ImagesIDs {
			choc, err := s.BmImageStorage.GetOne(kID)
			if err != nil {
				return &Response{}, err
			}
			model.Images = append(model.Images, &choc)
		}

		if model.CategoryID != "" {
			cate, err := s.BmCategoryStorage.GetOne(model.CategoryID)
			if err != nil {
				return &Response{}, err
			}
			model.Category = cate
		}

		result = append(result, *model)
	}

	return &Response{Res: result}, nil
}

// PaginatedFindAll can be used to load models in chunks
func (s BmSessioninfoResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []BmModel.Sessioninfo
		number, size, offset, limit string
	)

	numberQuery, ok := r.QueryParams["page[number]"]
	if ok {
		number = numberQuery[0]
	}
	sizeQuery, ok := r.QueryParams["page[size]"]
	if ok {
		size = sizeQuery[0]
	}
	offsetQuery, ok := r.QueryParams["page[offset]"]
	if ok {
		offset = offsetQuery[0]
	}
	limitQuery, ok := r.QueryParams["page[limit]"]
	if ok {
		limit = limitQuery[0]
	}

	if size != "" {
		sizeI, err := strconv.ParseInt(size, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		numberI, err := strconv.ParseInt(number, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		start := sizeI * (numberI - 1)
		for _, iter := range s.BmSessioninfoStorage.GetAll(r, int(start), int(sizeI)) {
			result = append(result, *iter)
		}

	} else {
		limitI, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		offsetI, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		for _, iter := range s.BmSessioninfoStorage.GetAll(r, int(offsetI), int(limitI)) {
			result = append(result, *iter)
		}
	}

	in := BmModel.Sessioninfo{}
	count := s.BmSessioninfoStorage.Count(r, in)

	return uint(count), &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the model with the given ID, otherwise an error
func (s BmSessioninfoResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmSessioninfoStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	model.Images = []*BmModel.Image{}
	for _, kID := range model.ImagesIDs {
		choc, err := s.BmImageStorage.GetOne(kID)
		if err != nil {
			return &Response{}, err
		}
		model.Images = append(model.Images, &choc)
	}
	return &Response{Res: model}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BmSessioninfoResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Sessioninfo)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := s.BmSessioninfoStorage.Insert(model)
	model.ID = id

	//TODO: 临时版本-在创建的同时加关系
	if model.CategoryID != "" {
		cate, err := s.BmCategoryStorage.GetOne(model.CategoryID)
		if err != nil {
			return &Response{}, err
		}
		model.Category = cate
	}

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s BmSessioninfoResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.BmSessioninfoStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the model
func (s BmSessioninfoResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Sessioninfo)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BmSessioninfoStorage.Update(model)
	return &Response{Res: model, Code: http.StatusNoContent}, err
}
