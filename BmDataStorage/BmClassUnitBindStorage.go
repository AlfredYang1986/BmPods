package BmDataStorage

import (
	"fmt"
	"github.com/alfredyang1986/BmPods/BmDaemons"
	"github.com/alfredyang1986/BmPods/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/manyminds/api2go"
)

// BmClassUnitBindStorage stores all of the tasty chocolate, needs to be injected into
// ClassUnitBind Resource. In the real world, you would use a database for that.
type BmClassUnitBindStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmClassUnitBindStorage) NewClassUnitBindStorage(args []BmDaemons.BmDaemon) *BmClassUnitBindStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmClassUnitBindStorage{mdb}
}

// GetAll of the chocolate
func (s BmClassUnitBindStorage) GetAll(r api2go.Request, skip int, take int) []*BmModel.ClassUnitBind {
	in := BmModel.ClassUnitBind{}
	var out []BmModel.ClassUnitBind
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.ClassUnitBind
		for i := 0; i < len(out); i++ {
			ptr := out[i]
			s.db.ResetIdWithId_(&ptr)
			tmp = append(tmp, &ptr)
		}
		return tmp
	} else {
		return nil //make(map[string]*BmModel.Student)
	}
}

// GetOne tasty chocolate
func (s BmClassUnitBindStorage) GetOne(id string) (BmModel.ClassUnitBind, error) {
	in := BmModel.ClassUnitBind{ID: id}
	out := BmModel.ClassUnitBind{ID: id}
	err := s.db.FindOne(&in, &out)
	/*if err == nil {

		if out.ClassID != "" {
			item, err := BmClassStorage{db: s.db}.GetOne(out.ClassID)
			if err != nil {
				return BmModel.ClassUnitBind{}, err
			}
			out.TeacherID = item.ID
		}
		if out.UnitID != "" {
			item, err := BmUnitStorage{db: s.db}.GetOne(out.UnitID)
			if err != nil {
				return BmModel.ClassUnitBind{}, err
			}
			out.Teacher = item
		}

		return out, nil
	}
	errMessage := fmt.Sprintf("ClassUnitBind for id %s not found", id)
	return BmModel.ClassUnitBind{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)*/
	return BmModel.ClassUnitBind{}, err
}

// Insert a fresh one
func (s *BmClassUnitBindStorage) Insert(c BmModel.ClassUnitBind) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmClassUnitBindStorage) Delete(id string) error {
	in := BmModel.ClassUnitBind{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("ClassUnitBind with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing chocolate
func (s *BmClassUnitBindStorage) Update(c BmModel.ClassUnitBind) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("ClassUnitBind with id does not exist")
	}

	return nil
}

func (s *BmClassUnitBindStorage) Count(req api2go.Request, c BmModel.ClassUnitBind) int {
	r, _ := s.db.Count(req, &c)
	return r
}