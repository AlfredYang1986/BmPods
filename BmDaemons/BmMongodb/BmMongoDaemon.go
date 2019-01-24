package BmMongodb

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/alfredyang1986/blackmirror/bmmate"
	"github.com/manyminds/api2go"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	BMJson    string = "json"
	BMJsonAPI string = "jsonapi"
	BMMongo   string = "bson"
)

type BmMongodb struct {
	Host     string
	Port     string
	Database string
}

type NoPtr struct {
}

func (m BmMongodb) NewMongoDBDaemon(args map[string]string) *BmMongodb {
	return &BmMongodb{
		Host:     args["host"],
		Port:     args["port"],
		Database: args["database"]}
}

func (m *BmMongodb) InsertBmObject(ptr BmModel.BmModelBase) (string, error) {
	session, err := mgo.Dial(m.Host + ":" + m.Port)
	if err != nil {
		return "", errors.New("dial db error")
	}
	defer session.Close()

	oid := m.GenerateModelId_(ptr)
	v := reflect.ValueOf(ptr).Elem()
	cn := v.Type().Name()
	c := session.DB(m.Database).C(cn)

	rst, err := Struct2map(v)
	rst["_id"] = oid
	err = c.Insert(rst)
	if err == nil {
		return m.ResetIdWithId_(ptr), nil
	}
	return "", err
}

func (m *BmMongodb) ExistsBmObject(ptr BmModel.BmModelBase, out BmModel.BmModelBase) (bool, error) {
	session, err := mgo.Dial(m.Host + ":" + m.Port)
	if err != nil {
		return false, errors.New("dial db error")
	}
	defer session.Close()

	oid := m.ResetId_WithId(ptr)
	v := reflect.ValueOf(ptr).Elem()
	cn := v.Type().Name()
	c := session.DB(m.Database).C(cn)

	err = c.Find(bson.M{"_id": oid}).One(out)
	if err != nil {
		fmt.Println(err)
	}

	m.ResetIdWithId_(out)
	return true, nil
}

func (m *BmMongodb) FindOne(ptr BmModel.BmModelBase, out BmModel.BmModelBase) error {
	session, err := mgo.Dial(m.Host + ":" + m.Port)
	if err != nil {
		return errors.New("dial db error")
	}
	defer session.Close()

	oid := m.ResetId_WithId(ptr)
	v := reflect.ValueOf(ptr).Elem()
	cn := v.Type().Name()
	c := session.DB(m.Database).C(cn)

	err = c.Find(bson.M{"_id": oid}).One(out)
	if err != nil {
		return errors.New("query error")
	}

	m.ResetIdWithId_(out)
	return nil
}

func (m *BmMongodb) FindMulti(r api2go.Request, ptr BmModel.BmModelBase, out interface{}, skip int, take int) error {
	session, err := mgo.Dial(m.Host + ":" + m.Port)
	if err != nil {
		return errors.New("dial db error")
	}
	defer session.Close()

	v := reflect.ValueOf(ptr).Elem()
	cn := v.Type().Name()
	c := session.DB(m.Database).C(cn)

	if skip < 0 && take < 0 {
		err = c.Find(ptr.GetConditionsBsonM(r.QueryParams)).Sort(Request2SortCondi(r)).All(out)
	} else {
		err = c.Find(ptr.GetConditionsBsonM(r.QueryParams)).Skip(skip).Limit(take).Sort(Request2SortCondi(r)).All(out)
	}

	if err != nil {
		return errors.New("error find multi")
	}

	return nil
}

func (m *BmMongodb) Delete(ptr BmModel.BmModelBase) error {
	session, err := mgo.Dial(m.Host + ":" + m.Port)
	if err != nil {
		return errors.New("dial db error")
	}
	defer session.Close()

	oid := m.ResetId_WithId(ptr)
	v := reflect.ValueOf(ptr).Elem()
	cn := v.Type().Name()
	c := session.DB(m.Database).C(cn)

	err = c.Remove(bson.M{"_id": oid})
	if err != nil {
		return errors.New("error delete by id")
	}
	return nil
}

func (m *BmMongodb) Update(ptr BmModel.BmModelBase) error {
	session, err := mgo.Dial(m.Host + ":" + m.Port)
	if err != nil {
		return errors.New("dial db error")
	}
	defer session.Close()

	oid := m.ResetId_WithId(ptr)
	v := reflect.ValueOf(ptr).Elem()
	cn := v.Type().Name()
	c := session.DB(m.Database).C(cn)

	rst, err := Struct2map(v)
	rst["_id"] = oid
	err = c.Update(bson.M{"_id": oid}, rst)
	if err != nil {
		return errors.New("error update by id")
	}
	return nil
}

func (m *BmMongodb) Count(r api2go.Request, ptr BmModel.BmModelBase) (int, error) {
	session, err := mgo.Dial(m.Host + ":" + m.Port)
	if err != nil {
		return 0, errors.New("dial db error")
	}
	defer session.Close()

	//m.ResetId_WithId(ptr)
	v := reflect.ValueOf(ptr).Elem()
	cn := v.Type().Name()
	c := session.DB(m.Database).C(cn)

	result, err := c.Find(ptr.GetConditionsBsonM(r.QueryParams)).Count()
	if err != nil {
		return 0, errors.New("error count by id")
	}
	return result, nil
}

func (m *BmMongodb) GenerateModelId_(ptr BmModel.BmModelBase) bson.ObjectId {
	f := reflect.ValueOf(ptr).Elem().FieldByName("Id_")
	v := bson.NewObjectId()
	f.Set(reflect.ValueOf(v))
	return v
}

func (m *BmMongodb) ResetIdWithId_(ptr BmModel.BmModelBase) string {
	f := reflect.ValueOf(ptr).Elem().FieldByName("Id_")
	t := f.Interface().(bson.ObjectId)
	fs := reflect.ValueOf(ptr).Elem().FieldByName("ID")
	fs.SetString(t.Hex())
	return t.Hex()
}

func (m *BmMongodb) ResetId_WithId(ptr BmModel.BmModelBase) bson.ObjectId {
	fs := reflect.ValueOf(ptr).Elem().FieldByName("ID")
	t := fs.Interface().(string)
	f := reflect.ValueOf(ptr).Elem().FieldByName("Id_")
	v := bson.ObjectIdHex(t)
	f.Set(reflect.ValueOf(v))
	return v
}

func AttrValue(v reflect.Value) (interface{}, error) {
	switch v.Kind() {
	case reflect.Invalid:
		return nil, nil
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return v.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return v.Float(), nil
	case reflect.String:
		return v.String(), nil
	case reflect.Array, reflect.Slice:
		var rst []interface{}
		for i := 0; i < v.Len(); i++ {
			tmp, _ := AttrValue(v.Index(i))
			rst = append(rst, tmp)
		}
		return rst, nil
	case reflect.Map:
		rst := make(map[string]interface{})
		for _, key := range v.MapKeys() {
			kv := v.MapIndex(key)
			tmp, _ := AttrValue(kv)
			rst[key.String()] = tmp
		}
		return rst, nil
	case reflect.Interface:
		if bmmate.IsStruct(v) {
			if reflect.TypeOf(v.Interface()).Kind() == reflect.String {
				return AttrValue(reflect.ValueOf(v.Interface()))
			} else {
				return AttrValue(reflect.ValueOf(v.Interface()))
			}
		} else {
			return AttrValue(reflect.ValueOf(v.Interface()))
		}
	}

	return NoPtr{}, errors.New("not implement")
}

func Struct2map(v reflect.Value) (map[string]interface{}, error) {
	rst := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {

		fieldInfo := v.Type().Field(i) // a.reflect.struct.field
		fieldValue := v.Field(i)
		tag := fieldInfo.Tag // a.reflect.tag

		var name string
		if tag.Get(BMMongo) != "" {
			name = tag.Get(BMMongo)
		} else if tag.Get(BMMongo) == "-" {
			continue
		} else {
			//name = strings.ToLower(fieldInfo.Name)
			continue
		}

		if name == "id" || name == "_id" {
			continue
		}

		tmp, _ := AttrValue(fieldValue)
		rst[name] = tmp
	}

	return rst, nil
}

// FindOneByCondition TODO 临时，稍后与老铁一起重写
func (m *BmMongodb) FindOneByCondition(ptr BmModel.BmModelBase, out BmModel.BmModelBase, cond bson.M) error {
	session, err := mgo.Dial(m.Host + ":" + m.Port)
	if err != nil {
		return errors.New("dial db error")
	}
	defer session.Close()

	v := reflect.ValueOf(ptr).Elem()
	cn := v.Type().Name()
	c := session.DB(m.Database).C(cn)

	err = c.Find(cond).One(out)
	if err != nil {
		return err
	}

	m.ResetIdWithId_(out)
	return nil
}

func Request2SortCondi(r api2go.Request) string {
	for k, v := range r.QueryParams {
		if k == "orderby" {
			return v[0]
		}
	}
	//如果不传orderby,暂时默认数据create-time的倒序
	return "-create-time"
}
