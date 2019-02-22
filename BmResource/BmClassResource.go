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

type BmClassResource struct {
	BmClassStorage          *BmDataStorage.BmClassStorage
	BmStudentStorage        *BmDataStorage.BmStudentStorage
	BmDutyStorage           *BmDataStorage.BmDutyStorage
	BmUnitStorage           *BmDataStorage.BmUnitStorage
	BmYardStorage           *BmDataStorage.BmYardStorage
	BmSessioninfoStorage    *BmDataStorage.BmSessioninfoStorage
	BmReservableitemStorage *BmDataStorage.BmReservableitemStorage
}

func (s BmClassResource) NewClassResource(args []BmDataStorage.BmStorage) BmClassResource {
	var us *BmDataStorage.BmClassStorage
	var ys *BmDataStorage.BmYardStorage
	var ss *BmDataStorage.BmSessioninfoStorage
	var cs *BmDataStorage.BmStudentStorage
	var ds *BmDataStorage.BmDutyStorage
	var ns *BmDataStorage.BmUnitStorage
	var rs *BmDataStorage.BmReservableitemStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmClassStorage" {
			us = arg.(*BmDataStorage.BmClassStorage)
		} else if tp.Name() == "BmStudentStorage" {
			cs = arg.(*BmDataStorage.BmStudentStorage)
		} else if tp.Name() == "BmDutyStorage" {
			ds = arg.(*BmDataStorage.BmDutyStorage)
		} else if tp.Name() == "BmUnitStorage" {
			ns = arg.(*BmDataStorage.BmUnitStorage)
		} else if tp.Name() == "BmYardStorage" {
			ys = arg.(*BmDataStorage.BmYardStorage)
		} else if tp.Name() == "BmSessioninfoStorage" {
			ss = arg.(*BmDataStorage.BmSessioninfoStorage)
		} else if tp.Name() == "BmReservableitemStorage" {
			rs = arg.(*BmDataStorage.BmReservableitemStorage)
		}
	}
	return BmClassResource{BmClassStorage: us, BmYardStorage: ys, BmSessioninfoStorage: ss, BmStudentStorage: cs, BmDutyStorage: ds, BmUnitStorage: ns, BmReservableitemStorage: rs}
}

// FindAll to satisfy api2go data source interface
func (s BmClassResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []BmModel.Class

	//查詢 reservable 下的 classes
	reservableitemsID, ok := r.QueryParams["reservableitemsID"]
	if ok {
		modelRootID := reservableitemsID[0]
		modelRoot, err := s.BmReservableitemStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		for _, modelID := range modelRoot.ClassesIDs {
			model, err := s.BmClassStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			err = s.ResetReferencedModel(&model)
			if err != nil {
				return &Response{}, err
			}
			result = append(result, model)
		}
		return &Response{Res: result}, nil
	}

	flag, fok := r.QueryParams["flag"]
	if fok {
		flagInt, err := strconv.Atoi(flag[0])
		if err != nil {
			return &Response{}, err
		}
		models := s.BmClassStorage.GetAll(r, -1, -1)
		for _, model := range models {
			err := s.ResetReferencedModel(model)
			if err != nil {
				return &Response{}, err
			}
			err = s.FilterClassByFlag(model, flagInt)
			if err == nil {
				result = append(result, *model)
			}
		}
		return &Response{Res: result}, nil
	}

	models := s.BmClassStorage.GetAll(r, -1, -1)
	for _, model := range models {
		err := s.ResetReferencedModel(model)
		if err != nil {
			return &Response{}, err
		}
		result = append(result, *model)
	}

	return &Response{Res: result}, nil
}

func (s BmClassResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []BmModel.Class
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

	reservableitemsID, ok := r.QueryParams["reservableitemsID"]
	if ok {
		modelRootID := reservableitemsID[0]
		modelRoot, err := s.BmReservableitemStorage.GetOne(modelRootID)
		if err != nil {
			return uint(0), &Response{}, err
		}
		count = len(modelRoot.ClassesIDs)
		if skip >= count {
			return uint(0), &Response{}, err
		}
		endIndex := skip + take
		if endIndex > count {
			endIndex = count
		}
		for _, modelID := range modelRoot.ClassesIDs[skip:endIndex] {
			model, err := s.BmClassStorage.GetOne(modelID)
			if err != nil {
				return uint(0), &Response{}, err
			}
			err = s.ResetReferencedModel(&model)
			if err != nil {
				return uint(0), &Response{}, err
			}
			result = append(result, model)
		}
		pages = int(math.Ceil(float64(count) / float64(take)))
		return uint(count), &Response{Res: result, QueryRes: "classes", TotalPage: pages}, nil
	}

	flag, fok := r.QueryParams["flag"]
	if fok {
		flagInt, err := strconv.Atoi(flag[0])
		if err != nil {
			return 0, &Response{}, err
		}
		models := s.BmClassStorage.GetAll(r, -1, -1)
		for _, model := range models {
			err := s.ResetReferencedModel(model)
			if err != nil {
				return 0, &Response{}, err
			}
			err = s.FilterClassByFlag(model, flagInt)
			if err == nil {
				result = append(result, *model)
			}
		}
		count = len(result)
		if skip >= count {
			return uint(0), &Response{}, err
		}
		endIndex := skip + take
		if endIndex > count {
			endIndex = count
		}
		pages = int(math.Ceil(float64(count) / float64(take)))
		return uint(count), &Response{Res: result[skip:endIndex], QueryRes: "classes", TotalPage: pages}, nil
	}

	for _, model := range s.BmClassStorage.GetAll(r, skip, take) {
		err := s.ResetReferencedModel(model)
		if err != nil {
			return 0, &Response{}, err
		}
		result = append(result, *model)
	}

	in := BmModel.Class{}
	count = s.BmClassStorage.Count(r, in)
	pages = int(math.Ceil(float64(count) / float64(take)))

	return uint(count), &Response{Res: result, QueryRes: "classes", TotalPage: pages}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the modelRoot with the given ID, otherwise an error
func (s BmClassResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmClassStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}
	err = s.ResetReferencedModel(&model)
	if len(model.UnitsIDs) != 0 {
		model.ReSetCourseCount()
	}
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	return &Response{Res: model}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BmClassResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Class)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	model.CreateTime = float64(time.Now().UnixNano() / 1e6)
	id := s.BmClassStorage.Insert(model)
	model.ID = id

	//TODO: 临时版本-在创建的同时加关系
	if model.YardID != "" {
		yard, err := s.BmYardStorage.GetOne(model.YardID)
		if err != nil {
			return &Response{}, err
		}
		model.Yard = yard
	}
	if model.SessioninfoID != "" {
		item, err := s.BmSessioninfoStorage.GetOne(model.SessioninfoID)
		if err != nil {
			return &Response{}, err
		}
		model.Sessioninfo = item
	}

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s BmClassResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.BmClassStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the modelRoot
func (s BmClassResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Class)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BmClassStorage.Update(model)
	return &Response{Res: model, Code: http.StatusNoContent}, err
}

func (s BmClassResource) ResetReferencedModel(model *BmModel.Class) error {
	model.Students = []*BmModel.Student{}
	for _, tmpID := range model.StudentsIDs {
		choc, err := s.BmStudentStorage.GetOne(tmpID)
		if err != nil {
			return err
		}
		model.Students = append(model.Students, &choc)
	}
	model.Duties = []*BmModel.Duty{}
	for _, tmpID := range model.DutiesIDs {
		choc, err := s.BmDutyStorage.GetOne(tmpID)
		if err != nil {
			return err
		}
		model.Duties = append(model.Duties, &choc)
	}
	model.Units = []*BmModel.Unit{}
	for _, tmpID := range model.UnitsIDs {
		choc, err := s.BmUnitStorage.GetOne(tmpID)
		if err != nil {
			return err
		}
		model.Units = append(model.Units, &choc)
	}

	if model.YardID != "" {
		yard, err := s.BmYardStorage.GetOne(model.YardID)
		if err != nil {
			return err
		}
		model.Yard = yard
	}
	if model.SessioninfoID != "" {
		item, err := s.BmSessioninfoStorage.GetOne(model.SessioninfoID)
		if err != nil {
			return err
		}
		model.Sessioninfo = item
	}
	return nil
}

func (s BmClassResource) FilterClassByFlag(model *BmModel.Class, flag int) error {
	switch flag {
	case 0:
		if len(model.UnitsIDs) != 0 {
			model.ReSetCourseCount()
		}
		return nil
	case -1:
		if len(model.UnitsIDs) == 0 {
			return nil
		}
	case 1:
		if len(model.UnitsIDs) != 0 {
			var us BmModel.Units
			us = model.Units
			us.SortByEndDate(false)
			now := float64(time.Now().UnixNano() / 1e6)
			if us[0].EndDate > now {
				model.ReSetCourseCount()
				return nil
			}
		}
	case 2:
		if len(model.UnitsIDs) != 0 {
			var us BmModel.Units
			us = model.Units
			us.SortByEndDate(false)
			now := float64(time.Now().UnixNano() / 1e6)
			if us[0].EndDate <= now {
				model.ReSetCourseCount()
				return nil
			}
		}
	}
	return errors.New("not found")
}
