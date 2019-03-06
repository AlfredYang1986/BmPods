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
	BmDutyResource           *BmDutyResource
	BmStudentResource        *BmStudentResource
	BmReservableitemResource *BmReservableitemResource
}

func (s BmClassResource) NewClassResource(args []BmDataStorage.BmStorage) *BmClassResource {
	var us *BmDataStorage.BmClassStorage
	var ys *BmDataStorage.BmYardStorage
	var ss *BmDataStorage.BmSessioninfoStorage
	var sr *BmStudentResource
	var rr *BmReservableitemResource
	var dr *BmDutyResource
	//var bs *BmDataStorage.BmClassUnitBindStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmClassStorage" {
			us = arg.(*BmDataStorage.BmClassStorage)
		} else if tp.Name() == "BmStudentResource" {
			sr = arg.(*BmStudentResource)
		} else if tp.Name() == "BmSessioninfoStorage" {
			ss = arg.(*BmDataStorage.BmSessioninfoStorage)
		} else if tp.Name() == "BmDutyResource" {
			dr = arg.(*BmDutyResource)
		} else if tp.Name() == "BmYardStorage" {
			ys = arg.(*BmDataStorage.BmYardStorage)
		} else if tp.Name() == "BmReservableitemResource" {
			rr = arg.(interface{}).(*BmReservableitemResource)
		}
	}
	return &BmClassResource{BmClassStorage: us, BmYardStorage: ys, BmSessioninfoStorage: ss, BmStudentResource: sr, BmDutyResource: dr, BmReservableitemResource: rr}
}

// FindAll to satisfy api2go data source interface
func (s BmClassResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []BmModel.Class

	//查詢 reservable 下的 classes
	//_, ok := r.QueryParams["reservableitem-id"]
	//if ok {
	//	modelBinds := s.BmBindReservableClassStorage.GetAll(r)
	//	for _, modelBind := range modelBinds {
	//		model, err := s.BmClassStorage.GetOne(modelBind.ClassId)
	//		if err != nil {
	//			return &Response{}, err
	//		}
	//		err = s.ResetReferencedModel(&model)
	//		if err != nil {
	//			return &Response{}, err
	//		}
	//		if model.NotExist == 0 {
	//			result = append(result, model)
	//		}
	//	}
	//	return &Response{Res: result}, nil
	//}

	// 写这个的人必须弹鸡鸡，无视数据库的所有规则，强制遍历所有自己filter
	// 那还要数据库做啥呢 。。。。
	//flag, fok := r.QueryParams["flag"]
	//if fok {
	//	flagInt, err := strconv.Atoi(flag[0])
	//	if err != nil {
	//		return &Response{}, err
	//	}
	//	models := s.BmClassStorage.GetAll(r, -1, -1)
	//	for _, model := range models {
	//		err := s.ResetReferencedModel(model)
	//		if err != nil {
	//			return &Response{}, err
	//		}
	//		err = s.FilterClassByFlag(model, flagInt)
	//		if err == nil && model.NotExist == 0 {
	//			result = append(result, *model)
	//		}
	//	}
	//	return &Response{Res: result}, nil
	//}

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
		//err := s.ResetReferencedModel(model, &r)
		//if err != nil {
		//	return &Response{}, err
		//}
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
		err := s.ResetReferencedModel(model, &r)
		if err != nil {
			return 0, &Response{}, err
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

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the modelRoot with the given ID, otherwise an error
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

// Create method to satisfy `api2go.DataSource` interface
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

// Delete to satisfy `api2go.DataSource` interface
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
	}
	if model.Execute==1{
		panic("只允许停课，不允许删除")
	}else{
		panic("不允许删除")
	}

	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the modelRoot
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
	stuRes, err := s.BmStudentResource.FindAll(*r)
	studs := stuRes.Result().([]BmModel.Student)
	for i, _ := range studs {
		model.Students = append(model.Students, &studs[i])
	}

	model.Duties = []*BmModel.Duty{}
	r.QueryParams["dutiesids"] = model.DutiesIDs
	response, err := s.BmDutyResource.FindAll(*r)
	if err != nil {
		return err
	}
	items := response.Result()
	duties := items.([]BmModel.Duty)
	for i, _ := range duties {
		model.Duties = append(model.Duties, &duties[i])
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
		//item, err := s.BmReservableitemStorage.GetOne(model.ReservableID)
		response, err := s.BmReservableitemResource.FindOne(model.ReservableID, api2go.Request{})
		item := response.Result()
		if err != nil {
			return err
		}
		tmp := item.(BmModel.Reservableitem)
		model.Reservableitem = &tmp
	}
	return nil
}
func (s BmClassResource) FilterClassByFlag(model *BmModel.Class, flag int) error {
	/*switch flag {
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
	}*/
	return errors.New("not found")
}
