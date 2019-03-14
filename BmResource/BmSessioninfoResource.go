package BmResource

import (
	"errors"
	//"fmt"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"math"
	"net/http"
	"reflect"
	"strconv"
)

type BmSessioninfoResource struct {
	BmImageStorage          *BmDataStorage.BmImageStorage
	BmSessioninfoStorage    *BmDataStorage.BmSessioninfoStorage
	BmCategoryStorage       *BmDataStorage.BmCategoryStorage
	BmClassStorage          *BmDataStorage.BmClassStorage
	BmReservableitemStorage *BmDataStorage.BmReservableitemStorage
}

func (s BmSessioninfoResource) NewSessioninfoResource(args []BmDataStorage.BmStorage) BmSessioninfoResource {
	var us *BmDataStorage.BmSessioninfoStorage
	var ts *BmDataStorage.BmCategoryStorage
	var is *BmDataStorage.BmImageStorage
	var cs *BmDataStorage.BmClassStorage
	var rs *BmDataStorage.BmReservableitemStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmSessioninfoStorage" {
			us = arg.(*BmDataStorage.BmSessioninfoStorage)
		} else if tp.Name() == "BmImageStorage" {
			is = arg.(*BmDataStorage.BmImageStorage)
		} else if tp.Name() == "BmCategoryStorage" {
			ts = arg.(*BmDataStorage.BmCategoryStorage)
		} else if tp.Name() == "BmReservableitemStorage" {
			rs = arg.(*BmDataStorage.BmReservableitemStorage)
		} else if tp.Name() == "BmClassStorage" {
			cs = arg.(*BmDataStorage.BmClassStorage)
		}
	}
	return BmSessioninfoResource{
		BmSessioninfoStorage:    us,
		BmImageStorage:          is,
		BmCategoryStorage:       ts,
		BmReservableitemStorage: rs,
		BmClassStorage:          cs,
	}
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
			err = s.ResetReferencedModel(&model, &r)
			if err != nil {
				return &Response{}, err
			}

			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	classesID, ok := r.QueryParams["classesID"]
	if ok {
		modelRootID := classesID[0]
		modelRoot, err := s.BmClassStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.SessioninfoID
		if modelID != "" {
			model, err := s.BmSessioninfoStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}

			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	models := s.BmSessioninfoStorage.GetAll(r, -1, -1)

	return &Response{Res: models}, nil
}

// PaginatedFindAll can be used to load models in chunks
func (s BmSessioninfoResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		number, size, offset, limit string
		skip, take, pages           int
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

		skip = int(start)
		take = int(sizeI)

	} else {
		limitI, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		offsetI, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return 0, &Response{}, err
		}

		skip = int(offsetI)
		take = int(limitI)
	}

	models := s.BmSessioninfoStorage.GetAll(r, skip, take)

	in := BmModel.Sessioninfo{}
	count := s.BmSessioninfoStorage.Count(r, in)
	pages = int(math.Ceil(float64(count) / float64(take)))
	return uint(count), &Response{Res: models, QueryRes: "reservableitems", TotalPage: pages, TotalCount: count}, nil
}

func (s BmSessioninfoResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmSessioninfoStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}
	err = s.ResetReferencedModel(&model, &r)
	if err != nil {
		return &Response{}, err
	}
	return &Response{Res: model}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BmSessioninfoResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Sessioninfo)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	imagesIDs := []string{}
	for _, img := range model.ImagesIDs {
		if img != "" {
			imagesIDs = append(imagesIDs, img)
		}
	}
	model.ImagesIDs = imagesIDs

	id := s.BmSessioninfoStorage.Insert(model)
	model.ID = id

	//TODO: 临时版本-在创建的同时加关系
	if model.CategoryID != "" {
		cate, err := s.BmCategoryStorage.GetOne(model.CategoryID)
		if err != nil {
			return &Response{}, err
		}
		model.Category = &cate
	}

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

func (s BmSessioninfoResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	panic("sessioninfo不可被删除")
}

func (s BmSessioninfoResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Sessioninfo)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	imagesIDs := []string{}
	for _, img := range model.ImagesIDs {
		if img != "" {
			imagesIDs = append(imagesIDs, img)
		}
	}
	model.ImagesIDs = imagesIDs

	err := s.BmSessioninfoStorage.Update(model)
	return &Response{Res: model, Code: http.StatusNoContent}, err
}

func (s BmSessioninfoResource) ResetReferencedModel(model *BmModel.Sessioninfo, r *api2go.Request) error {

	model.Images = []*BmModel.Image{}
	r.QueryParams["imageids"] = model.ImagesIDs
	imageids := s.BmImageStorage.GetAll(*r)
	for i, _ := range imageids {
		model.Images = append(model.Images, &imageids[i])
	}
	if model.CategoryID != "" {
		cate, err := s.BmCategoryStorage.GetOne(model.CategoryID)
		if err != nil {
			return err
		}
		model.Category = &cate
	}

	return nil
}
