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

type BmDutyResource struct {
	BmDutyStorage    *BmDataStorage.BmDutyStorage
	BmTeacherStorage *BmDataStorage.BmTeacherStorage
}

func (s BmDutyResource) NewDutyResource(args []BmDataStorage.BmStorage) *BmDutyResource {
	var ds *BmDataStorage.BmDutyStorage
	var ts *BmDataStorage.BmTeacherStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmDutyStorage" {
			ds = arg.(*BmDataStorage.BmDutyStorage)
		} else if tp.Name() == "BmTeacherStorage" {
			ts = arg.(*BmDataStorage.BmTeacherStorage)
		}
	}
	return &BmDutyResource{BmDutyStorage: ds, BmTeacherStorage: ts}
}

// FindAll to satisfy api2go data source interface
func (s BmDutyResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []BmModel.Duty
	duties := s.BmDutyStorage.GetAll(r, -1, -1)

	for _, duty := range duties {
		// get all sweets for the duty
		if duty.TeacherID != "" {
			r, err := s.BmTeacherStorage.GetOne(duty.TeacherID)
			if err != nil {
				return &Response{}, errors.New("error")
			}
			duty.Teacher = &r
		}
		result = append(result, *duty)
	}

	return &Response{Res: result}, nil
}

// PaginatedFindAll can be used to load users in chunks
func (s BmDutyResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []BmModel.Duty
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
		for _, iter := range s.BmDutyStorage.GetAll(r, int(start), int(sizeI)) {
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

		for _, iter := range s.BmDutyStorage.GetAll(r, int(offsetI), int(limitI)) {
			result = append(result, *iter)
		}
	}

	in := BmModel.Duty{}
	count := s.BmDutyStorage.Count(r, in)

	return uint(count), &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the user with the given ID, otherwise an error
func (s BmDutyResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmDutyStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	if model.TeacherID != "" {
		r, err := s.BmTeacherStorage.GetOne(model.TeacherID)
		if err != nil {
			return &Response{}, errors.New("error")
		}
		model.Teacher = &r
	}

	return &Response{Res: model}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BmDutyResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Duty)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := s.BmDutyStorage.Insert(model)
	model.ID = id

	if model.TeacherID != "" {
		r, err := s.BmTeacherStorage.GetOne(model.TeacherID)
		if err != nil {
			return &Response{}, errors.New("error")
		}
		model.Teacher = &r
	}

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s BmDutyResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.BmDutyStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the user
func (s BmDutyResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	user, ok := obj.(BmModel.Duty)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BmDutyStorage.Update(user)
	return &Response{Res: user, Code: http.StatusNoContent}, err
}
