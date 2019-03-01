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

type BmBrandResource struct {
	BmBrandStorage    *BmDataStorage.BmBrandStorage
	BmImageStorage    *BmDataStorage.BmImageStorage
	BmCategoryStorage *BmDataStorage.BmCategoryStorage
}

func (s BmBrandResource) NewBrandResource(args []BmDataStorage.BmStorage) BmBrandResource {
	var bs *BmDataStorage.BmBrandStorage
	var is *BmDataStorage.BmImageStorage
	var cs *BmDataStorage.BmCategoryStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmBrandStorage" {
			bs = arg.(*BmDataStorage.BmBrandStorage)
		} else if tp.Name() == "BmImageStorage" {
			is = arg.(*BmDataStorage.BmImageStorage)
		} else if tp.Name() == "BmCategoryStorage" {
			cs = arg.(*BmDataStorage.BmCategoryStorage)
		}
	}
	return BmBrandResource{BmBrandStorage: bs, BmImageStorage: is, BmCategoryStorage: cs}
}

// FindAll to satisfy api2go data source interface
func (s BmBrandResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []BmModel.Brand
	models := s.BmBrandStorage.GetAll(r, -1, -1)

	for _, model := range models {
		// get all sweets for the model
		model.Imgs = []*BmModel.Image{}
		r.QueryParams["imageids"]=model.ImagesIDs
		imageids:=s.BmImageStorage.GetAll(r)
		for _,image:= range imageids {	
			model.Imgs = append(model.Imgs, &image)
		}
/*
		model.Imgs = []*BmModel.Image{}
		for _, kID := range model.ImagesIDs {
			choc, err := s.BmImageStorage.GetOne(kID)
			if err != nil {
				return &Response{}, err
			}
			model.Imgs = append(model.Imgs, &choc)
		}
*/
		if model.CategoryID != "" {
			cat, err := s.BmCategoryStorage.GetOne(model.CategoryID)
			if err != nil {
				return &Response{}, err
			}
			model.Cat = cat
		}

		result = append(result, *model)
	}

	return &Response{Res: result}, nil
}

// PaginatedFindAll can be used to load models in chunks
func (s BmBrandResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []BmModel.Brand
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
		for _, iter := range s.BmBrandStorage.GetAll(r, int(start), int(sizeI)) {
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

		for _, iter := range s.BmBrandStorage.GetAll(r, int(offsetI), int(limitI)) {
			result = append(result, *iter)
		}
	}

	in := BmModel.Brand{}
	count := s.BmBrandStorage.Count(r, in)

	return uint(count), &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the model with the given ID, otherwise an error
func (s BmBrandResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmBrandStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	model.Imgs = []*BmModel.Image{}
		r.QueryParams["imageids"]=model.ImagesIDs
		imageids:=s.BmImageStorage.GetAll(r)
		for _,image:= range imageids {	
			model.Imgs = append(model.Imgs, &image)
		}

	if model.CategoryID != "" {
		cat, err := s.BmCategoryStorage.GetOne(model.CategoryID)
		if err != nil {
			return &Response{}, err
		}
		model.Cat = cat
	}

	return &Response{Res: model}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BmBrandResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Brand)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := s.BmBrandStorage.Insert(model)
	model.ID = id

	//TODO: 临时版本-在创建的同时加关系
	if model.CategoryID != "" {
		cat, err := s.BmCategoryStorage.GetOne(model.CategoryID)
		if err != nil {
			return &Response{}, err
		}
		model.Cat = cat
	}

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s BmBrandResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.BmBrandStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the model
func (s BmBrandResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Brand)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BmBrandStorage.Update(model)
	return &Response{Res: model, Code: http.StatusNoContent}, err
}
