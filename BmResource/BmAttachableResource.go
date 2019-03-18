package BmResource

import (
	"errors"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
)

type BmAttachableResource struct {
	BmTransactionStorage *BmDataStorage.BmTransactionStorage
	BmAttachableStorage *BmDataStorage.BmAttachableStorage
	BmStudentStorage    *BmDataStorage.BmStudentStorage
}

func (c BmAttachableResource) NewAttachableResource(args []BmDataStorage.BmStorage) BmAttachableResource {
	var ts *BmDataStorage.BmTransactionStorage
	var as *BmDataStorage.BmAttachableStorage
	var ss *BmDataStorage.BmStudentStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmAttachableStorage" {
			as = arg.(*BmDataStorage.BmAttachableStorage)
		} else if tp.Name() == "BmStudentStorage" {
			ss = arg.(*BmDataStorage.BmStudentStorage)
		} else if tp.Name() == "BmTransactionStorage" {
			ts = arg.(*BmDataStorage.BmTransactionStorage)
		}
	}
	return BmAttachableResource{BmAttachableStorage: as, BmStudentStorage: ss, BmTransactionStorage: ts}
}

func (c BmAttachableResource) FindAll(r api2go.Request) (api2go.Responder, error) {

	transactionsID, ok := r.QueryParams["transactionsID"]
	if ok {
		modelRootID := transactionsID[0]
		modelRoot, err := c.BmTransactionStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		r.QueryParams["attachableids"] = modelRoot.AttachablesIDs
		results := c.BmAttachableStorage.GetAll(r, -1, -1)
		for i, _ := range results {
			c.ResetReferencedModel(&results[i], &r)
		}
		return &Response{Res: results}, nil
	}

	results := c.BmAttachableStorage.GetAll(r, -1, -1)
	return &Response{Res: results}, nil
}

func (c BmAttachableResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.BmAttachableStorage.GetOne(ID)
	if err != nil {
		return &Response{}, err
	}
	err = c.ResetReferencedModel(&res, &r)
	if err != nil {
		return &Response{}, err
	}
	return &Response{Res: res}, err
}

func (c BmAttachableResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Attachable)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := c.BmAttachableStorage.Insert(choc)
	choc.ID = id
	return &Response{Res: choc, Code: http.StatusCreated}, nil
}

func (c BmAttachableResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	panic("不可删除")
}

func (c BmAttachableResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Attachable)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := c.BmAttachableStorage.Update(choc)
	return &Response{Res: choc, Code: http.StatusNoContent}, err
}

func (c BmAttachableResource) ResetReferencedModel(model *BmModel.Attachable, r *api2go.Request) error {

	if model.StudentID != "" {
		leaf, err := c.BmStudentStorage.GetOne(model.StudentID)
		if err != nil {
			return err
		}
		model.Student = &leaf
	}

	return nil
}
