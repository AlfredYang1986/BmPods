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
	BmClassStorage           *BmDataStorage.BmClassStorage
	BmYardStorage            *BmDataStorage.BmYardStorage
	BmSessioninfoStorage     *BmDataStorage.BmSessioninfoStorage
	BmDutyStorage            *BmDataStorage.BmDutyStorage
	BmUnitStorage            *BmDataStorage.BmUnitStorage
	BmStudentStorage         *BmDataStorage.BmStudentStorage
	BmReservableitemStorage  *BmDataStorage.BmReservableitemStorage
}

func (s BmClassResource) NewClassResource(args []BmDataStorage.BmStorage) *BmClassResource {
	var cs *BmDataStorage.BmClassStorage
	var ys *BmDataStorage.BmYardStorage
	var ss *BmDataStorage.BmSessioninfoStorage
	var sr *BmDataStorage.BmStudentStorage
	var rs *BmDataStorage.BmReservableitemStorage
	var ds *BmDataStorage.BmDutyStorage
	var us *BmDataStorage.BmUnitStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmClassStorage" {
			cs = arg.(*BmDataStorage.BmClassStorage)
		} else if tp.Name() == "BmStudentStorage" {
			sr = arg.(*BmDataStorage.BmStudentStorage)
		} else if tp.Name() == "BmSessioninfoStorage" {
			ss = arg.(*BmDataStorage.BmSessioninfoStorage)
		} else if tp.Name() == "BmDutyStorage" {
			ds = arg.(*BmDataStorage.BmDutyStorage)
		} else if tp.Name() == "BmYardStorage" {
			ys = arg.(*BmDataStorage.BmYardStorage)
		} else if tp.Name() == "BmReservableitemStorage" {
			rs = arg.(*BmDataStorage.BmReservableitemStorage)
		} else if tp.Name() == "BmUnitStorage" {
			us = arg.(*BmDataStorage.BmUnitStorage)
		}
	}
	return &BmClassResource{BmClassStorage: cs, BmYardStorage: ys, BmSessioninfoStorage: ss, BmStudentStorage: sr, BmDutyStorage: ds, BmReservableitemStorage: rs, BmUnitStorage: us}
}

func (s BmClassResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []BmModel.Class


	models := s.BmClassStorage.GetAll(r, -1, -1)
	for _, model := range models {
		now := float64(time.Now().UnixNano() / 1e6)
		if len(model.DutiesIDs) != 0 || len(model.StudentsIDs) != 0 || model.ReservableID != "" {
			model.Execute = 1
		}
		if now > model.StartDate && now <= model.EndDate {
			model.Execute = 2
		} else if now > model.EndDate {
			model.Execute = 3
		}
		if model.NotExist == 0 {
			result = append(result, *model)
		}
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

	for _, model := range s.BmClassStorage.GetAll(r, skip, take) {
		now := float64(time.Now().UnixNano() / 1e6)
		if len(model.DutiesIDs) != 0 || len(model.StudentsIDs) != 0 || model.ReservableID != "" {
			model.Execute = 1
		}
		if now > model.StartDate && now <= model.EndDate {
			model.Execute = 2
		} else if now > model.EndDate {
			model.Execute = 3
		}
		if model.NotExist == 0 {
			result = append(result, *model)
		}
	}

	in := BmModel.Class{}
	count = s.BmClassStorage.Count(r, in)
	pages = int(math.Ceil(float64(count) / float64(take)))

	return uint(count), &Response{Res: result, QueryRes: "classes", TotalPage: pages, TotalCount: count}, nil
}

func (s BmClassResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmClassStorage.GetOne(ID)
	now := float64(time.Now().UnixNano() / 1e6)
	if len(model.DutiesIDs) != 0 || len(model.StudentsIDs) != 0 || model.ReservableID != "" {
		model.Execute = 1
	}
	if now > model.StartDate && now <= model.EndDate {
		model.Execute = 2
	} else if now > model.EndDate {
		model.Execute = 3
	}
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}
	err = s.ResetReferencedModel(&model, &r)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}
	return &Response{Res: model}, nil
}

func (s BmClassResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Class)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	model.CreateTime = float64(time.Now().UnixNano() / 1e6)
	id := s.BmClassStorage.Insert(model)
	model.ID = id
	s.ResetReferencedModel(&model, &r)
	return &Response{Res: model, Code: http.StatusCreated}, nil
}

func (s BmClassResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmClassStorage.GetOne(id)
	now := float64(time.Now().UnixNano() / 1e6)
	if len(model.DutiesIDs) != 0 || len(model.StudentsIDs) != 0 || model.ReservableID != "" {
		model.Execute = 1
	}
	if now > model.StartDate && now <= model.EndDate {
		model.Execute = 2
	} else if now > model.EndDate {
		model.Execute = 3
	}
	if err != nil {
		return &Response{}, err
	}
	if model.Execute == 0 {
		s.BmClassStorage.Delete(id)
		r.QueryParams["class-id"] = []string{id}
		us := s.BmUnitStorage.GetAll(r, -1, -1)
		if len(us) != 0 {
			for _, u := range us {
				s.BmUnitStorage.Delete(u.ID)
			}
		}

	}
	if model.Execute==1{
		panic("只允许停课，不允许删除")
	}else{
		panic("不允许删除")
	}

	return &Response{Code: http.StatusNoContent}, err
}

func (s BmClassResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Class)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BmClassStorage.Update(model)
	s.ResetReferencedModel(&model, &r)
	return &Response{Res: model, Code: http.StatusNoContent}, err
}

func (s BmClassResource) ResetReferencedModel(model *BmModel.Class, r *api2go.Request) error {
	model.Students = []*BmModel.Student{}
	
	r.QueryParams["studentsids"] = model.StudentsIDs
	stuRes:= s.BmStudentStorage.GetAll(*r,-1,-1)
	for i, _ := range stuRes {
		model.Students = append(model.Students, stuRes[i])
	}

	model.Duties = []*BmModel.Duty{}
	r.QueryParams["dutiesids"] = model.DutiesIDs
	dutRes := s.BmDutyStorage.GetAll(*r,-1,-1)
	for i, _ := range dutRes {
		model.Duties = append(model.Duties, dutRes[i])
	}

	if model.YardID != "" {
		yard, err := s.BmYardStorage.GetOne(model.YardID)
		if err != nil {
			return err
		}
		model.Yard = &yard
	}
	if model.SessioninfoID != "" {
		sessioninfo, err := s.BmSessioninfoStorage.GetOne(model.SessioninfoID)
		if err != nil {
			return err
		}
		model.Sessioninfo = &sessioninfo
	}

	if model.ReservableID != "" {
		reservableitem, err := s.BmReservableitemStorage.GetOne(model.ReservableID)
		if err != nil {
			return err
		}
		model.Reservableitem = &reservableitem
	}
	return nil
}
func (s BmClassResource) FilterClassByFlag(model *BmModel.Class, flag int) error {
	return errors.New("not found")
}
