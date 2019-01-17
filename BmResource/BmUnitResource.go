package BmResource

import (
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"reflect"
	"github.com/manyminds/api2go"
	"github.com/alfredyang1986/BmPods/BmModel"
	"strconv"
	"net/http"
	"errors"
)

type BmUnitResource struct {
	BmUnitStorage *BmDataStorage.BmUnitStorage
	BmRoomStorage *BmDataStorage.BmRoomStorage
	BmTeacherStorage *BmDataStorage.BmTeacherStorage
}

func (s BmUnitResource) NewUnitResource(args []BmDataStorage.BmStorage) BmUnitResource {
	var us *BmDataStorage.BmUnitStorage
	var rs *BmDataStorage.BmRoomStorage
	var ts *BmDataStorage.BmTeacherStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmBmUnitStorage" {
			us = arg.(*BmDataStorage.BmUnitStorage)
		} else if tp.Name() == "BmRoomStorage" {
			rs = arg.(*BmDataStorage.BmRoomStorage)
		} else if tp.Name() == "BmTeacherStorage" {
			ts = arg.(*BmDataStorage.BmTeacherStorage)
		}
	}
	return BmUnitResource{BmUnitStorage: us, BmRoomStorage: rs, BmTeacherStorage: ts}
}

// FindAll to satisfy api2go data source interface
func (s BmUnitResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []BmModel.Unit
	users := s.BmUnitStorage.GetAll(-1, -1)

	for _, user := range users {
		// get all sweets for the user

		if user.RoomID != "" {
			r, err := s.BmRoomStorage.GetOne(user.RoomID)
			if err != nil {
				return &Response{}, errors.New("error")
			}
			user.Room = &r
		}

		if user.TeacherID != "" {
			r, err := s.BmTeacherStorage.GetOne(user.TeacherID)
			if err != nil {
				return &Response{}, errors.New("error")
			}
			user.Teacher = &r
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
		for _, iter := range s.BmUnitStorage.GetAll(int(start), int(sizeI)) {
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

		for _, iter := range s.BmUnitStorage.GetAll(int(offsetI), int(limitI)) {
			result = append(result, *iter)
		}
	}

	in := BmModel.Unit{}
	count := s.BmUnitStorage.Count(in)

	return uint(count), &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the user with the given ID, otherwise an error
func (s BmUnitResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	user, err := s.BmUnitStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	if user.RoomID != "" {
		r, err := s.BmRoomStorage.GetOne(user.RoomID)
		if err != nil {
			return &Response{}, errors.New("error")
		}
		user.Room = &r
	}

	if user.TeacherID != "" {
		r, err := s.BmTeacherStorage.GetOne(user.TeacherID)
		if err != nil {
			return &Response{}, errors.New("error")
		}
		user.Teacher = &r
	}

	return &Response{Res: user}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BmUnitResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	user, ok := obj.(BmModel.Unit)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := s.BmUnitStorage.Insert(user)
	user.ID = id

	return &Response{Res: user, Code: http.StatusCreated}, nil
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