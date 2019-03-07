package BmResource

import (
	"errors"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type BmTeacherResource struct {
	BmTeacherStorage *BmDataStorage.BmTeacherStorage
	BmDutyStorage    *BmDataStorage.BmDutyStorage
	BmUnitStorage    *BmDataStorage.BmUnitStorage
	BmStudentStorage    *BmDataStorage.BmStudentStorage
}

func (s BmTeacherResource) NewTeacherResource(args []BmDataStorage.BmStorage) BmTeacherResource {
	var ts *BmDataStorage.BmTeacherStorage
	var ds *BmDataStorage.BmDutyStorage
	var us *BmDataStorage.BmUnitStorage
	var ss *BmDataStorage.BmStudentStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmTeacherStorage" {
			ts = arg.(*BmDataStorage.BmTeacherStorage)
		} else if tp.Name() == "BmDutyStorage" {
			ds = arg.(*BmDataStorage.BmDutyStorage)
		} else if tp.Name() == "BmUnitStorage" {
			us = arg.(*BmDataStorage.BmUnitStorage)
		}
	}
	return BmTeacherResource{BmTeacherStorage: ts, BmDutyStorage: ds, BmUnitStorage: us, BmStudentStorage: ss}
}

// FindAll to satisfy api2go data source interface
func (s BmTeacherResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	result := []BmModel.Teacher{}

	studentsID, ok := r.QueryParams["studentsID"]
	if ok {
		modelRootID := studentsID[0]
		modelRoot, err := s.BmStudentStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.TeacherID
		if modelID != "" {
			model, err := s.BmTeacherStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	dutiesID, ok := r.QueryParams["dutiesID"]
	if ok {
		modelRootID := dutiesID[0]
		modelRoot, err := s.BmDutyStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.TeacherID
		if modelID != "" {
			model, err := s.BmTeacherStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	unitsID, ok := r.QueryParams["unitsID"]
	if ok {
		modelRootID := unitsID[0]
		modelRoot, err := s.BmUnitStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.TeacherID
		if modelID != "" {
			model, err := s.BmTeacherStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	result = s.BmTeacherStorage.GetAll(r, -1, -1)

	return &Response{Res: result}, nil
}

func (s BmTeacherResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []BmModel.Teacher
		number, size, offset, limit string
		skip, take, count           int
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
	result = s.BmTeacherStorage.GetAll(r, skip, take)
	in := BmModel.Teacher{}
	count = s.BmTeacherStorage.Count(r, in)

	return uint(count), &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the user with the given ID, otherwise an error
func (s BmTeacherResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmTeacherStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	return &Response{Res: model}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BmTeacherResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Teacher)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	model.CreateTime = float64(time.Now().UnixNano() / 1e6)
	id := s.BmTeacherStorage.Insert(model)
	model.ID = id

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s BmTeacherResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmTeacherStorage.GetOne(id)
	if err != nil {
		return &Response{}, err
	}
	model.Archive = 1.0
	err = s.BmTeacherStorage.Update(model)
	if err != nil {
		return &Response{}, err
	}
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the user
func (s BmTeacherResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Teacher)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BmTeacherStorage.Update(model)
	return &Response{Res: model, Code: http.StatusNoContent}, err
}
