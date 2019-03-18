package BmResource

import (
	"errors"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"gopkg.in/mgo.v2/bson"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"time"
)

type BmTransactionResource struct {
	BmTransactionStorage *BmDataStorage.BmTransactionStorage
	BmAttachableStorage  *BmDataStorage.BmAttachableStorage
	BmApplicantStorage   *BmDataStorage.BmApplicantStorage
	BmTeacherStorage     *BmDataStorage.BmTeacherStorage
	BmStudentStorage     *BmDataStorage.BmStudentStorage
}

func (c BmTransactionResource) NewTransactionResource(args []BmDataStorage.BmStorage) BmTransactionResource {
	var trs *BmDataStorage.BmTransactionStorage
	var ats *BmDataStorage.BmAttachableStorage
	var aps *BmDataStorage.BmApplicantStorage
	var tes *BmDataStorage.BmTeacherStorage
	var sts *BmDataStorage.BmStudentStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmTransactionStorage" {
			trs = arg.(*BmDataStorage.BmTransactionStorage)
		} else if tp.Name() == "BmAttachableStorage" {
			ats = arg.(*BmDataStorage.BmAttachableStorage)
		} else if tp.Name() == "BmApplicantStorage" {
			aps = arg.(*BmDataStorage.BmApplicantStorage)
		} else if tp.Name() == "BmTeacherStorage" {
			tes = arg.(*BmDataStorage.BmTeacherStorage)
		} else if tp.Name() == "BmStudentStorage" {
			sts = arg.(*BmDataStorage.BmStudentStorage)
		}
	}
	return BmTransactionResource{BmTransactionStorage: trs, BmAttachableStorage: ats, BmApplicantStorage: aps, BmTeacherStorage: tes, BmStudentStorage: sts}
}

func (c BmTransactionResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	transactions := c.BmTransactionStorage.GetAll(r, -1, -1)
	return &Response{Res: transactions}, nil
}

// PaginatedFindAll can be used to load models in chunks
func (s BmTransactionResource) PaginatedFindAll(r api2go.Request) (uint, api2go.Responder, error) {
	var (
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

	results := s.BmTransactionStorage.GetAll(r, skip, take)
	in := BmModel.Transaction{}
	count = s.BmTransactionStorage.Count(r, in)
	pages = int(math.Ceil(float64(count) / float64(take)))
	return uint(count), &Response{Res: results, QueryRes: "transactions", TotalPage: pages, TotalCount: count}, nil
}

func (c BmTransactionResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.BmTransactionStorage.GetOne(ID)
	if err != nil {
		return &Response{}, err
	}
	err = c.ResetReferencedModel(&res, &r)
	if err != nil {
		return &Response{}, err
	}
	return &Response{Res: res}, err
}

func (c BmTransactionResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Transaction)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	choc.Id_ = bson.NewObjectId()
	choc.ID = choc.Id_.Hex()
	choc.OrderId = "BMID-" + choc.ID
	choc.CreateTime = float64(time.Now().UnixNano() / 1e6)
	id := c.BmTransactionStorage.Insert(choc)
	choc.ID = id
	return &Response{Res: choc, Code: http.StatusCreated}, nil
}

func (c BmTransactionResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	panic("不可删除")
}

func (c BmTransactionResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Transaction)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := c.BmTransactionStorage.Update(choc)
	return &Response{Res: choc, Code: http.StatusNoContent}, err
}

func (c BmTransactionResource) ResetReferencedModel(model *BmModel.Transaction, r *api2go.Request) error {

	model.Attachables = []*BmModel.Attachable{}
	r.QueryParams["attachableids"] = model.AttachablesIDs
	attachables := c.BmAttachableStorage.GetAll(*r, -1, -1)
	for i, v := range attachables {
		if v.StudentID != "" {
			leaf, err := c.BmStudentStorage.GetOne(v.StudentID)
			if err != nil {
				return err
			}
			attachables[i].Student = &leaf
		}
		model.Attachables = append(model.Attachables, &attachables[i])
	}
	if model.ApplicantID != "" {
		leaf, err := c.BmApplicantStorage.GetOne(model.ApplicantID)
		if err != nil {
			return err
		}
		model.Applicant = &leaf
	}
	if model.TeacherID != "" {
		leaf, err := c.BmTeacherStorage.GetOne(model.TeacherID)
		if err != nil {
			return err
		}
		model.Teacher = &leaf
	}

	return nil
}
