package BmResource

import (
	"errors"
	"github.com/alfredyang1986/BmPods/BmDataStorage"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
)

type BmGuardianResource struct {
	BmGuardianStorage *BmDataStorage.BmGuardianStorage
	BmStudentStorage  *BmDataStorage.BmStudentStorage
}

func (c BmGuardianResource) NewGuardianResource(args []BmDataStorage.BmStorage) BmGuardianResource {
	var cs *BmDataStorage.BmGuardianStorage
	var ss *BmDataStorage.BmStudentStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmGuardianStorage" {
			cs = arg.(*BmDataStorage.BmGuardianStorage)
		}
		if tp.Name() == "BmStudentStorage" {
			ss = arg.(*BmDataStorage.BmStudentStorage)
		}
	}
	return BmGuardianResource{BmGuardianStorage: cs, BmStudentStorage: ss}
}

// FindAll guardians
func (c BmGuardianResource) FindAll(r api2go.Request) (api2go.Responder, error) {

	result := []BmModel.Guardian{}
	studentsID, ok := r.QueryParams["studentsID"]
	if ok {
		modelRootID := studentsID[0]
		modelRoot, err := c.BmStudentStorage.GetOne(modelRootID)
		if err != nil {
			return &Response{}, err
		}
		if len(modelRoot.GuardiansIDs) != 0 {
			r.QueryParams["guardiansids"] = modelRoot.GuardiansIDs
			models := c.BmGuardianStorage.GetAll(r)
			return &Response{Res: models}, nil
		}

		return &Response{Res: result}, nil
	}

	result = c.BmGuardianStorage.GetAll(r)
	return &Response{Res: result}, nil
}

// FindOne choc
func (c BmGuardianResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.BmGuardianStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new choc
func (c BmGuardianResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Guardian)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := c.BmGuardianStorage.Insert(choc)
	choc.ID = id
	return &Response{Res: choc, Code: http.StatusCreated}, nil
}

// Delete a choc :(
func (c BmGuardianResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.BmGuardianStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a choc
func (c BmGuardianResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Guardian)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := c.BmGuardianStorage.Update(choc)
	return &Response{Res: choc, Code: http.StatusNoContent}, err
}
