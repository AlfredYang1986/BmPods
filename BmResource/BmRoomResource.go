package BmResource

import (
	"errors"
	"math"
	"net/http"
	"reflect"
	"strconv"

	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
)

type BmRoomResource struct {
	BmRoomStorage *BmDataStorage.BmRoomStorage
	BmYardStorage *BmDataStorage.BmYardStorage
}

func (c BmRoomResource) NewRoomResource(args []BmDataStorage.BmStorage) BmRoomResource {
	var us *BmDataStorage.BmYardStorage
	var cs *BmDataStorage.BmRoomStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmYardStorage" {
			us = arg.(*BmDataStorage.BmYardStorage)
		} else if tp.Name() == "BmRoomStorage" {
			cs = arg.(*BmDataStorage.BmRoomStorage)
		}
	}
	return BmRoomResource{BmYardStorage: us, BmRoomStorage: cs}
}

func (c BmRoomResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	result := []BmModel.Room{}
	yardsID, ok := r.QueryParams["yardsID"]
	if ok {
		modelRootID := yardsID[0]
		modelRoot, err := c.BmYardStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		for _, modelID := range modelRoot.RoomsIDs {
			model, err := c.BmRoomStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			result = append(result, model)
		}

		return &Response{Res: result}, nil
	}
	result = c.BmRoomStorage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

func (s BmRoomResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []BmModel.Room
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
			return uint(0), &Response{}, err
		}

		numberI, err := strconv.ParseInt(number, 10, 64)
		if err != nil {
			return uint(0), &Response{}, err
		}

		start := sizeI * (numberI - 1)

		skip = int(start)
		take = int(sizeI)
	} else {
		limitI, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return uint(0), &Response{}, err
		}

		offsetI, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return uint(0), &Response{}, err
		}

		skip = int(offsetI)
		take = int(limitI)
	}

	yardsID, ok := r.QueryParams["yardsID"]
	if ok {
		modelRootID := yardsID[0]
		modelRoot, err := s.BmYardStorage.GetOne(modelRootID)
		if err != nil {
			return uint(0), &Response{}, err
		}
		count = len(modelRoot.RoomsIDs)
		if skip >= count {
			return uint(0), &Response{}, err
		}
		endIndex := skip + take
		if endIndex >= count {
			endIndex = count
		}
		for _, modelID := range modelRoot.RoomsIDs[skip:endIndex] {
			model, err := s.BmRoomStorage.GetOne(modelID)
			if err != nil {
				return uint(0), &Response{}, err
			}
			result = append(result, model)
		}
		pages = int(math.Ceil(float64(count) / float64(take)))
		return uint(count), &Response{Res: result, QueryRes: "rooms", TotalPage: pages}, nil
	}

	result = s.BmRoomStorage.GetAll(r, skip, take)
	in := BmModel.Room{}
	count = s.BmRoomStorage.Count(r, in)
	pages = int(math.Ceil(float64(count) / float64(take)))
	return uint(count), &Response{Res: result, QueryRes: "rooms", TotalPage: pages}, nil
}

// FindOne choc
func (c BmRoomResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.BmRoomStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new choc
func (c BmRoomResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Room)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := c.BmRoomStorage.Insert(choc)
	choc.ID = id
	return &Response{Res: choc, Code: http.StatusCreated}, nil
}

// Delete a choc :(
func (c BmRoomResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	panic("不可以删除教室")
}

// Update a choc
func (c BmRoomResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Room)
	var err error
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}
	if choc.Archive==0&&choc.IsUnit==0{
		err = c.BmRoomStorage.Update(choc)
	}else if choc.Archive==1{
		panic("该房间未开放，不允许修改")
	}
	return &Response{Res: choc, Code: http.StatusNoContent}, err
}
