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

type BmStudentResource struct {
	BmStudentStorage    *BmDataStorage.BmStudentStorage
	BmKidStorage        *BmDataStorage.BmKidStorage
	BmTeacherStorage    *BmDataStorage.BmTeacherStorage
	BmGuardianStorage	*BmDataStorage.BmGuardianStorage
	BmClassStorage   	*BmDataStorage.BmClassStorage
	BmApplicantStorage *BmDataStorage.BmApplicantStorage
}

func (s BmStudentResource) NewStudentResource(args []BmDataStorage.BmStorage) *BmStudentResource {
	var ss *BmDataStorage.BmStudentStorage
	var ks *BmDataStorage.BmKidStorage
	var gs *BmDataStorage.BmGuardianStorage
	var ts *BmDataStorage.BmTeacherStorage
	var cs *BmDataStorage.BmClassStorage
	var as *BmDataStorage.BmApplicantStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmStudentStorage" {
			ss = arg.(*BmDataStorage.BmStudentStorage)
		} else if tp.Name() == "BmKidStorage" {
			ks = arg.(*BmDataStorage.BmKidStorage)
		} else if tp.Name() == "BmGuardianStorage" {
			gs = arg.(*BmDataStorage.BmGuardianStorage)
		} else if tp.Name() == "BmTeacherStorage" {
			ts = arg.(*BmDataStorage.BmTeacherStorage)
		} else if tp.Name() == "BmClassStorage" {
			cs = arg.(*BmDataStorage.BmClassStorage)
		} else if tp.Name() == "BmApplicantStorage" {
			as = arg.(*BmDataStorage.BmApplicantStorage)
		}
	}
	return &BmStudentResource{BmStudentStorage: ss, BmKidStorage: ks, BmGuardianStorage: gs, BmTeacherStorage: ts, BmClassStorage: cs,BmApplicantStorage: as}
}

// FindAll to satisfy api2go data source interface
func (s BmStudentResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []BmModel.Student
	var kidsids []string
	contact, ok := r.QueryParams["contact"]
	r.QueryParams["regi-phone"]=contact
	if ok {
		applients := s.BmApplicantStorage.GetAll(r,-1,-1)
		for _,applient:=range applients{
			r.QueryParams["applicant-id"]=[]string{applient.ID}
			kids := s.BmKidStorage.GetAll(r)
			for _,kid := range kids{	
				kidsids=append(kidsids,kid.ID)
			}
			r.QueryParams["kidids"]=kidsids
			students := s.BmStudentStorage.GetAll(r,-1,-1)
			for _,student:=range students{
				result = append(result,*student)
			}
		}
		return &Response{Res: result}, nil
	}
	//查詢 class 下的 students
	classesID, ok := r.QueryParams["classesID"]
	if ok {
		modelRootID := classesID[0]
		modelRoot, err := s.BmClassStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		for _, modelID := range modelRoot.StudentsIDs {
			model, err := s.BmStudentStorage.GetOne(modelID)
			if err != nil {
				return &Response{}, err
			}
			err = s.ResetReferencedModel(&model,&r)
			if err != nil {
				return &Response{}, err
			}

			result = append(result, model)
		}
		return &Response{Res: result}, nil
	}

	models := s.BmStudentStorage.GetAll(r, -1, -1)
	return &Response{Res: models}, nil
}

// PaginatedFindAll can be used to load users in chunks
func (s BmStudentResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
		result                      []BmModel.Student
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

	var kidsids []string
	contact, ok := r.QueryParams["contact"]
	r.QueryParams["regi-phone"]=contact
	if ok {
		applients := s.BmApplicantStorage.GetAll(r,-1,-1)
		for _,applient:=range applients{
			r.QueryParams["applicant-id"]=[]string{applient.ID}
			kids := s.BmKidStorage.GetAll(r)
			for _,kid := range kids{	
				kidsids=append(kidsids,kid.ID)
			}
			r.QueryParams["kidids"]=kidsids
			students := s.BmStudentStorage.GetAll(r,skip,take)
			for _,student:=range students{
				result = append(result,*student)
			}
		}
		count = len(result)
		pages = int(math.Ceil(float64(count) / float64(take)))
		return uint(count), &Response{Res: result, QueryRes: "students", TotalPage: pages, TotalCount: count}, nil
	}
	//查詢class下的students
	classesID, ok := r.QueryParams["classesID"]
	if ok {
		modelRootID := classesID[0]
		modelRoot, err := s.BmClassStorage.GetOne(modelRootID)
		if err != nil {
			return uint(0), &Response{}, err
		}
		count = len(modelRoot.StudentsIDs)
		if skip >= count {
			return uint(0), &Response{}, err
		}
		endIndex := skip + take
		if endIndex >= count {
			endIndex = count
		}
		for _, modelID := range modelRoot.StudentsIDs[skip:endIndex] {
			model, err := s.BmStudentStorage.GetOne(modelID)
			if err != nil {
				return uint(0), &Response{}, err
			}
			err = s.ResetReferencedModel(&model,&r)
			if err != nil {
				return uint(0), &Response{}, err
			}

			result = append(result, model)
		}
		pages = int(math.Ceil(float64(count) / float64(take)))
		return uint(count), &Response{Res: result, QueryRes: "students", TotalPage: pages, TotalCount: count}, nil
	}

	models:=s.BmStudentStorage.GetAll(r, skip, take) 
		
	in := BmModel.Student{}
	count = s.BmStudentStorage.Count(r, in)
	pages = int(math.Ceil(float64(count) / float64(take)))
	return uint(count), &Response{Res: models, QueryRes: "students", TotalPage: pages}, nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the user with the given ID, otherwise an error
func (s BmStudentResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	model, err := s.BmStudentStorage.GetOne(ID)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}
	err = s.ResetReferencedModel(&model,&r)
	if err != nil {
		return &Response{}, api2go.NewHTTPError(err, err.Error(), http.StatusNotFound)
	}

	return &Response{Res: model}, nil
}

// Create method to satisfy `api2go.DataSource` interface
func (s BmStudentResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	model, ok := obj.(BmModel.Student)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	model.CreateTime = float64(time.Now().UnixNano() / 1e6)
	id := s.BmStudentStorage.Insert(model)
	model.ID = id

	//TODO: 临时版本-在创建的同时加关系
	if model.KidID != "" {
		k, err := s.BmKidStorage.GetOne(model.KidID)
		if err != nil {
			return &Response{}, err
		}
		model.Kid = &k
	}
	if model.TeacherID != "" {
		k, err := s.BmTeacherStorage.GetOne(model.TeacherID)
		if err != nil {
			return &Response{}, err
		}
		model.Teacher = &k
	}

	return &Response{Res: model, Code: http.StatusCreated}, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s BmStudentResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := s.BmKidStorage.Delete(id)
	return &Response{Code: http.StatusNoContent}, err
}

//Update stores all changes on the user
func (s BmStudentResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	user, ok := obj.(BmModel.Student)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := s.BmStudentStorage.Update(user)
	return &Response{Res: user, Code: http.StatusNoContent}, err
}

func (s BmStudentResource) ResetReferencedModel(model *BmModel.Student,r *api2go.Request) error {

	model.Guardians = []*BmModel.Guardian{}
	r.QueryParams["guardiansids"]=model.GuardiansIDs
	guardians:=s.BmGuardianStorage.GetAll(*r)
	for _,guardian:= range guardians {	
		model.Guardians = append(model.Guardians, &guardian)
	}
/*			
	for _, chocolateID := range model.GuardiansIDs {
		choc, err := s.BmGuardianStorage.GetOne(chocolateID)
		if err != nil {
			return err
		}
		model.Guardians = append(model.Guardians, &choc)
	}
*/
	if model.KidID != "" {
		k, err := s.BmKidStorage.GetOne(model.KidID)
		if err != nil {
			return err
		}
		model.Kid = &k
	}

	if model.TeacherID != "" {
		k, err := s.BmTeacherStorage.GetOne(model.TeacherID)
		if err != nil {
			return err
		}
		model.Teacher = &k
	}

	return nil
}
