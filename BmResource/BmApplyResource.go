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

type BmApplyResource struct {
	BmKidStorage       *BmDataStorage.BmKidStorage
	BmApplyStorage     *BmDataStorage.BmApplyStorage
	BmApplicantStorage *BmDataStorage.BmApplicantStorage
}

func (s BmApplyResource) NewApplyResource(args []BmDataStorage.BmStorage) BmApplyResource {
	var us *BmDataStorage.BmApplyStorage
	var ts *BmDataStorage.BmApplicantStorage
	var cs *BmDataStorage.BmKidStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmApplyStorage" {
			us = arg.(*BmDataStorage.BmApplyStorage)
		} else if tp.Name() == "BmKidStorage" {
			cs = arg.(*BmDataStorage.BmKidStorage)
		} else if tp.Name() == "BmApplicantStorage" {
			ts = arg.(*BmDataStorage.BmApplicantStorage)
		}
	}
	return BmApplyResource{BmApplyStorage: us, BmKidStorage: cs, BmApplicantStorage: ts}
}

// FindAll to satisfy api2go data source interface
func (s BmApplyResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []BmModel.Apply
	models := s.BmApplyStorage.GetAll(r, -1, -1)

	for _, model := range models {
		// get all sweets for the model
		model.Kids = []*BmModel.Kid{}
		for _, kID := range model.KidsIDs {
			choc, err := s.BmKidStorage.GetOne(kID)
			if err != nil {
				return &Response{}, err
			}
			model.Kids = append(model.Kids, &choc)
		}

		if model.ApplicantID != "" {
			applicant, err := s.BmApplicantStorage.GetOne(model.ApplicantID)
			if err != nil {
				return &Response{}, err
			}
			model.Applicant = applicant
		}

		result = append(result, *model)
	}

	return &Response{Res: result}, nil
}

// PaginatedFindAll can be used to load models in chunks
func (s BmApplyResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []BmModel.Apply
		number, size, offset, limit string
		skip, take, count, pages    int
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

	for _, model := range s.BmApplyStorage.GetAll(r, skip, take) {

		model.Kids = []*BmModel.Kid{}
		for _, kID := range model.KidsIDs {
			choc, err := s.BmKidStorage.GetOne(kID)
			if err != nil {
				return 0, &Response{}, err
			}
			model.Kids = append(model.Kids, &choc)
		}

		if model.ApplicantID != "" {
			applicant, err := s.BmApplicantStorage.GetOne(model.ApplicantID)
			if err != nil {
				return 0, &Response{}, err
			}
			model.Applicant = applicant
		}

		result = append(result, *model)
	}

	in := BmModel.Apply{}
	count = s.BmApplyStorage.Count(r, in)
	pages = int(math.Ceil(float64(count) / float64(take)))
	return uint(count), &Response{Res: result, QueryRes: "applies", TotalPage: pages, TotalCount: count}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the model with the given ID, otherwise an error
func (s BmApplyResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmApplyStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	if model.ApplicantID != "" {
		applicant, err := s.BmApplicantStorage.GetOne(model.ApplicantID)
		if err != nil {
			return &Response{}, err
		}
		model.Applicant = applicant
	}

	model.Kids = []*BmModel.Kid{}
	for _, kID := range model.KidsIDs {
		choc, err := s.BmKidStorage.GetOne(kID)
		if err != nil {
			return &Response{}, err
		}
		model.Kids = append(model.Kids, &choc)
	}
	return &Response{Res: model}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BmApplyResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Apply)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	model.CreateTime = float64(time.Now().UnixNano() / 1e6)
	id := s.BmApplyStorage.Insert(model)
	model.ID = id

	//TODO: 临时版本-在创建的同时加关系
	if model.ApplicantID != "" {
		applicant, err := s.BmApplicantStorage.GetOne(model.ApplicantID)
		if err != nil {
			return &Response{}, err
		}
		model.Applicant = applicant
	}

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s BmApplyResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.BmApplyStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the model
func (s BmApplyResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Apply)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BmApplyStorage.Update(model)
	return &Response{Res: model, Code: http.StatusNoContent}, err
}
