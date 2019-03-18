package BmResource

import (
	"errors"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type BmReservableitemResource struct {
	BmReservableitemStorage *BmDataStorage.BmReservableitemStorage
	BmSessioninfoStorage    *BmDataStorage.BmSessioninfoStorage
	BmCategoryStorage       *BmDataStorage.BmCategoryStorage
	BmImageStorage          *BmDataStorage.BmImageStorage
	BmClassStorage          *BmDataStorage.BmClassStorage
}

func (s BmReservableitemResource) NewReservableitemResource(args []BmDataStorage.BmStorage) *BmReservableitemResource {
	var us *BmDataStorage.BmReservableitemStorage
	var ts *BmDataStorage.BmSessioninfoStorage
	var gs *BmDataStorage.BmCategoryStorage
	var is *BmDataStorage.BmImageStorage
	var cs *BmDataStorage.BmClassStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmReservableitemStorage" {
			us = arg.(*BmDataStorage.BmReservableitemStorage)
		} else if tp.Name() == "BmSessioninfoStorage" {
			ts = arg.(*BmDataStorage.BmSessioninfoStorage)
		} else if tp.Name() == "BmCategoryStorage" {
			gs = arg.(*BmDataStorage.BmCategoryStorage)
		} else if tp.Name() == "BmImageStorage" {
			is = arg.(*BmDataStorage.BmImageStorage)
		} else if tp.Name() == "BmClassStorage" {
			cs = arg.(*BmDataStorage.BmClassStorage)
		}
	}
	return &BmReservableitemResource{BmReservableitemStorage: us, BmSessioninfoStorage: ts, BmCategoryStorage: gs, BmImageStorage: is, BmClassStorage: cs}
}

// FindAll to satisfy api2go data source interface
func (s BmReservableitemResource) FindAll(r api2go.Request) (api2go.Responder, error) {

	var result []BmModel.Reservableitem
	_, ok := r.QueryParams["sub-title"]
	if ok {
		sessioninfos := s.BmSessioninfoStorage.GetAll(r,-1,-1)
		for _,sessioninfo:=range sessioninfos{
			r.QueryParams["sessioninfo-id"]=[]string{sessioninfo.ID}
			reservableitems := s.BmReservableitemStorage.GetAll(r,-1,-1)		
			for _,reservableitem:=range reservableitems{
				result = append(result,*reservableitem)
			}
			return &Response{Res: result}, nil
		}
	}
	
	_, titleok := r.QueryParams["title"]
	if titleok {
		sessioninfos := s.BmSessioninfoStorage.GetAll(r,-1,-1)
		for _,sessioninfo:=range sessioninfos{
			r.QueryParams["sessioninfo-id"]=[]string{sessioninfo.ID}
			reservableitems := s.BmReservableitemStorage.GetAll(r,-1,-1)		
			for _,reservableitem:=range reservableitems{
				result = append(result,*reservableitem)
			}
			return &Response{Res: result}, nil
		}
	}

	classesID, ok := r.QueryParams["classesID"]
	if ok {
		modelRootID := classesID[0]
		modelRoot, err := s.BmClassStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.ReservableID
		if modelID != "" {
			model, err := s.BmReservableitemStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			//result = append(result, model)
			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	models := s.BmReservableitemStorage.GetAll(r, -1, -1)
	return &Response{Res: models}, nil
}

// PaginatedFindAll can be used to load models in chunks
func (s BmReservableitemResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []BmModel.Reservableitem
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
	for _, model := range s.BmReservableitemStorage.GetAll(r, skip, take) {
		now := float64(time.Now().UnixNano() / 1e6)
		if now <= model.StartDate {
			model.Execute = 0
		} else if now > model.StartDate && now <= model.EndDate {
			model.Execute = 2
		} else {
			model.Execute = 1
		}
		result = append(result, *model)
	}

	in := BmModel.Reservableitem{}
	count := s.BmReservableitemStorage.Count(r, in)
	pages = int(math.Ceil(float64(count) / float64(take)))
	return uint(count), &Response{Res: result, QueryRes: "reservableitems", TotalPage: pages, TotalCount: count}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the model with the given ID, otherwise an error
func (s BmReservableitemResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmReservableitemStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}
	now := float64(time.Now().UnixNano() / 1e6)
	if now <= model.StartDate {
		model.Execute = 0
	} else if now > model.StartDate && now <= model.EndDate {
		model.Execute = 2
	} else {
		model.Execute = 1
	}
	if model.SessioninfoID != "" {
		sessioninfo, err := s.BmSessioninfoStorage.GetOne(model.SessioninfoID)
		if err != nil {
			return &Response{}, err
		}

		sessioninfo.Images = []*BmModel.Image{}
		r.QueryParams["imageids"] = sessioninfo.ImagesIDs
		imageids := s.BmImageStorage.GetAll(r)
		for i, _ := range imageids {
			sessioninfo.Images = append(sessioninfo.Images, &imageids[i])
		}
		if sessioninfo.CategoryID != "" {
			cate, err := s.BmCategoryStorage.GetOne(sessioninfo.CategoryID)
			if err != nil {
				return &Response{}, err
			}
			sessioninfo.Category = &cate
		}

		model.Sessioninfo = &sessioninfo
	}
	return &Response{Res: model}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BmReservableitemResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Reservableitem)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	model.CreateTime = float64(time.Now().UnixNano() / 1e6)
	id := s.BmReservableitemStorage.Insert(model)
	model.ID = id

	//TODO: 临时版本-在创建的同时加关系
	if model.SessioninfoID != "" {
		sessioninfo, err := s.BmSessioninfoStorage.GetOne(model.SessioninfoID)
		if err != nil {
			return &Response{}, err
		}
		model.Sessioninfo = &sessioninfo
	}

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s BmReservableitemResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	panic("不可删除")
}

//Update stores all changes on the model
func (s BmReservableitemResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Reservableitem)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BmReservableitemStorage.Update(model)
	return &Response{Res: model, Code: http.StatusNoContent}, err
}
