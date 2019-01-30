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

type BmUnitResource struct {
	BmUnitStorage    *BmDataStorage.BmUnitStorage
	BmRoomStorage    *BmDataStorage.BmRoomStorage
	BmTeacherStorage *BmDataStorage.BmTeacherStorage
	BmClassStorage   *BmDataStorage.BmClassStorage
}

func (s BmUnitResource) NewUnitResource(args []BmDataStorage.BmStorage) BmUnitResource {
	var us *BmDataStorage.BmUnitStorage
	var rs *BmDataStorage.BmRoomStorage
	var ts *BmDataStorage.BmTeacherStorage
	var cs *BmDataStorage.BmClassStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmUnitStorage" {
			us = arg.(*BmDataStorage.BmUnitStorage)
		} else if tp.Name() == "BmRoomStorage" {
			rs = arg.(*BmDataStorage.BmRoomStorage)
		} else if tp.Name() == "BmTeacherStorage" {
			ts = arg.(*BmDataStorage.BmTeacherStorage)
		} else if tp.Name() == "BmClassStorage" {
			cs = arg.(*BmDataStorage.BmClassStorage)
		}
	}
	return BmUnitResource{BmUnitStorage: us, BmRoomStorage: rs, BmTeacherStorage: ts, BmClassStorage: cs}
}

// FindAll to satisfy api2go data source interface
func (s BmUnitResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []BmModel.Unit
	users := s.BmUnitStorage.GetAll(r, -1, -1)

	for _, user := range users {
		// get all sweets for the user

		if user.RoomID != "" {
			r, err := s.BmRoomStorage.GetOne(user.RoomID)
			if err != nil {
				return &Response{}, errors.New("error")
			}
			user.Room = r
		}

		if user.TeacherID != "" {
			r, err := s.BmTeacherStorage.GetOne(user.TeacherID)
			if err != nil {
				return &Response{}, errors.New("error")
			}
			user.Teacher = r
		}
		if user.ClassID != "" {
			r, err := s.BmClassStorage.GetOne(user.ClassID)
			if err != nil {
				return &Response{}, errors.New("error")
			}
			user.Class = r
		}

		result = append(result, *user)
	}

	return &Response{Res: result}, nil
}

// PaginatedFindAll can be used to load users in chunks
func (s BmUnitResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []BmModel.Unit
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
		for _, iter := range s.BmUnitStorage.GetAll(r, int(start), int(sizeI)) {
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

		for _, iter := range s.BmUnitStorage.GetAll(r, int(offsetI), int(limitI)) {
			result = append(result, *iter)
		}
	}

	in := BmModel.Unit{}
	count := s.BmUnitStorage.Count(r, in)

	return uint(count), &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the user with the given ID, otherwise an error
func (s BmUnitResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmUnitStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	if model.RoomID != "" {
		r, err := s.BmRoomStorage.GetOne(model.RoomID)
		if err != nil {
			return &Response{}, errors.New("error")
		}
		model.Room = r
	}

	if model.TeacherID != "" {
		r, err := s.BmTeacherStorage.GetOne(model.TeacherID)
		if err != nil {
			return &Response{}, errors.New("error")
		}
		model.Teacher = r
	}
	if model.ClassID != "" {
		r, err := s.BmClassStorage.GetOne(model.ClassID)
		if err != nil {
			return &Response{}, errors.New("error")
		}
		model.Class = r
	}

	return &Response{Res: model}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BmUnitResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Unit)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := s.BmUnitStorage.Insert(model)
	model.ID = id

	if model.RoomID != "" {
		r, err := s.BmRoomStorage.GetOne(model.RoomID)
		if err != nil {
			return &Response{}, errors.New("error")
		}
		model.Room = r
	}

	if model.TeacherID != "" {
		r, err := s.BmTeacherStorage.GetOne(model.TeacherID)
		if err != nil {
			return &Response{}, errors.New("error")
		}
		model.Teacher = r
	}
	if model.ClassID != "" {
		r, err := s.BmClassStorage.GetOne(model.ClassID)
		if err != nil {
			return &Response{}, errors.New("error")
		}
		model.Class = r
	}

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s BmUnitResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.BmUnitStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the user
func (s BmUnitResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	user, ok := obj.(BmModel.Unit)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BmUnitStorage.Update(user)
	return &Response{Res: user, Code: http.StatusNoContent}, err
}
