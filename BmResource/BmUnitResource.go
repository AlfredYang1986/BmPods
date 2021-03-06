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
		}else if tp.Name() == "BmClassStorage" {
			cs = arg.(*BmDataStorage.BmClassStorage)
		} 
	}
	return BmUnitResource{BmUnitStorage: us, BmRoomStorage: rs, BmTeacherStorage: ts, BmClassStorage: cs}
}

func (s BmUnitResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []BmModel.Unit

	models := s.BmUnitStorage.GetAll(r, -1, -1)
	
	for _, model := range models {
		now := float64(time.Now().UnixNano() / 1e6)
		if now <= model.StartDate {
			model.Execute=0
		}else if now > model.StartDate && now <= model.EndDate{
			model.Execute=2
		}else{
			model.Execute=1
		}
		result = append(result, *model)
	}

	return &Response{Res: result}, nil
}

func (s BmUnitResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		number, size, offset, limit string
		skip, take, count   int
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
	result:= s.BmUnitStorage.GetAll(r,skip, take) 		
	for _, model := range result {
		now := float64(time.Now().UnixNano() / 1e6)
		if now <= model.StartDate {
			model.Execute=0
		}else if now > model.StartDate && now <= model.EndDate{
			model.Execute=2
		}else{
			model.Execute=1
		}
		result = append(result, model)
	}

	in := BmModel.Unit{}
	count = s.BmUnitStorage.Count(r, in)

	return uint(count), &Response{Res: result}, nil
}

func (s BmUnitResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmUnitStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}
	now := float64(time.Now().UnixNano() / 1e6)
	if now <= model.StartDate {
		model.Execute=0
	}else if now > model.StartDate && now <= model.EndDate{
		model.Execute=2
	}else{
		model.Execute=1
	}
	if model.RoomID != "" {
		r, err := s.BmRoomStorage.GetOne(model.RoomID)
		if err != nil {
			return &Response{}, err
		}
		model.Room = &r
	}

	if model.TeacherID != "" {
		r, err := s.BmTeacherStorage.GetOne(model.TeacherID)
		if err != nil {
			return &Response{}, err
		}
		model.Teacher = &r
	}

	if model.ClassID != "" {
		r, err := s.BmClassStorage.GetOne(model.ClassID)
		if err != nil {
			return &Response{}, err
		}
		model.Class = &r
	}


	return &Response{Res: model}, nil
}

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
			return &Response{}, err
		}
		model.Room = &r
	}

	if model.TeacherID != "" {
		r, err := s.BmTeacherStorage.GetOne(model.TeacherID)
		if err != nil {
			return &Response{}, err
		}
		model.Teacher = &r
	}

	if model.ClassID != "" {
		r, err := s.BmClassStorage.GetOne(model.ClassID)
		if err != nil {
			return &Response{}, err
		}
		model.Class = &r
	}

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

func (s BmUnitResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmUnitStorage.GetOne(id)
	if err != nil {
		return &Response{}, err
	}
	now := float64(time.Now().UnixNano() / 1e6)
	if now <= model.StartDate {
		model.Execute=0
	}else if now > model.StartDate && now <= model.EndDate{
		model.Execute=2
	}else{
		model.Execute=1
	}
	if model.Execute==0{
		model.Archive = 1.0
		err = s.BmUnitStorage.Update(model)
		if err != nil {
			return &Response{}, err
		}
	}
	if model.Execute==1{
		panic("已结束，不可删除")
	}

	return &Response{Code: http.StatusNoContent}, err
}

func (s BmUnitResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Unit)
	var err error 
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}
	room, err:=s.BmRoomStorage.GetOne(model.RoomID)
	if err != nil {
		return &Response{}, err
	}
	if room.IsUnit==0{
		room.IsUnit=1
		err = s.BmRoomStorage.Update(room)
		if err != nil {
			return &Response{}, err
		}
	}
	now := float64(time.Now().UnixNano() / 1e6)
	if now <= model.StartDate {
		model.Execute=0
	}else if now > model.StartDate && now <= model.EndDate{
		model.Execute=2
	}else{
		model.Execute=1
	}

	if model.Execute==1{
		panic("已结束，不可编辑")
	} else {
		err = s.BmUnitStorage.Update(model)
		if err != nil {
			return &Response{}, err
		}
	}

	return &Response{Res: model, Code: http.StatusNoContent}, err
}
