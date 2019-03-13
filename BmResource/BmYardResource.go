package BmResource

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
)

type BmYardResource struct {
	BmYardStorage  *BmDataStorage.BmYardStorage
	BmImageStorage *BmDataStorage.BmImageStorage
	BmRoomStorage  *BmDataStorage.BmRoomStorage
	BmClassStorage  *BmDataStorage.BmClassStorage
}

func (s BmYardResource) NewYardResource(args []BmDataStorage.BmStorage) BmYardResource {
	var ys *BmDataStorage.BmYardStorage
	var is *BmDataStorage.BmImageStorage
	var rs *BmDataStorage.BmRoomStorage
	var cs *BmDataStorage.BmClassStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmYardStorage" {
			ys = arg.(*BmDataStorage.BmYardStorage)
		} else if tp.Name() == "BmImageStorage" {
			is = arg.(*BmDataStorage.BmImageStorage)
		} else if tp.Name() == "BmRoomStorage" {
			rs = arg.(*BmDataStorage.BmRoomStorage)
		} else if tp.Name() == "BmClassStorage" {
			cs = arg.(*BmDataStorage.BmClassStorage)
		}
	}
	return BmYardResource{BmYardStorage: ys, BmImageStorage: is, BmRoomStorage: rs, BmClassStorage: cs}
}

// FindAll to satisfy api2go data source interface
func (s BmYardResource) FindAll(r api2go.Request) (api2go.Responder, error) {

	classesID, ok := r.QueryParams["classesID"]
	if ok {
		modelRootID := classesID[0]
		modelRoot, err := s.BmClassStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		modelID := modelRoot.YardID
		if modelID != "" {
			model, err := s.BmYardStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			return &Response{Res: model}, nil
		} else {
			return &Response{}, err
		}
	}

	models := s.BmYardStorage.GetAll(r, -1, -1)
	return &Response{Res: models}, nil
}

// PaginatedFindAll can be used to load models in chunks
func (s BmYardResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []BmModel.Yard
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
		for _, iter := range s.BmYardStorage.GetAll(r, int(start), int(sizeI)) {
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

		for _, iter := range s.BmYardStorage.GetAll(r, int(offsetI), int(limitI)) {
			result = append(result, *iter)
		}
	}

	in := BmModel.Yard{}
	count := s.BmYardStorage.Count(r, in)

	return uint(count), &Response{Res: result}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the model with the given ID, otherwise an error
func (s BmYardResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmYardStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	model.Images = []*BmModel.Image{}
	r.QueryParams["imageids"] = model.ImagesIDs
	images := s.BmImageStorage.GetAll(r)
	for i, _ := range images {
		model.Images = append(model.Images, &images[i])
	}

	model.Rooms = []*BmModel.Room{}
	r.QueryParams["roomids"] = model.ImagesIDs
	rooms := s.BmRoomStorage.GetAll(r, -1, -1)
	for i, _ := range rooms {
		model.Rooms = append(model.Rooms, &rooms[i])
	}

	return &Response{Res: model}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BmYardResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Yard)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := s.BmYardStorage.Insert(model)
	model.ID = id

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s BmYardResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.BmYardStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the model
func (s BmYardResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Yard)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BmYardStorage.Update(model)
	return &Response{Res: model, Code: http.StatusNoContent}, err
}
