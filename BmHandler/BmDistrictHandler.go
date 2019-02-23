package BmHandler

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmPods/BmDaemons/BmMongodb"
	"github.com/julienschmidt/httprouter"
)

type DistrictHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
}

func (h DistrictHandler) NewDistrictHandler(args ...interface{}) DistrictHandler {
	var m *BmMongodb.BmMongodb
	var hm string
	var md string
	var ag []string
	//sts := args[0].([]BmDaemons.BmDaemon)
	for i, arg := range args {
		if i == 0 {
			sts := arg.([]BmDaemons.BmDaemon)
			for _, dm := range sts {
				tp := reflect.ValueOf(dm).Interface()
				tm := reflect.ValueOf(tp).Elem().Type()
				if tm.Name() == "BmMongodb" {
					m = dm.(*BmMongodb.BmMongodb)
				}
			}
		} else if i == 1 {
			md = arg.(string)
		} else if i == 2 {
			hm = arg.(string)
		} else if i == 3 {
			lst := arg.([]string)
			for _, str := range lst {
				ag = append(ag, str)
			}
		} else {
		}
	}
	return DistrictHandler{Method: md, HttpMethod: hm, Args: ag, db: m}
}

//TODO: load files
func (h DistrictHandler) AllDistricts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	data := []string{
		"密云区",
		"朝阳区",
		"丰台区",
		"石景山区",
		"海淀区",
		"门头沟区",
		"房山区",
		"通州区",
		"顺义区",
		"昌平区",
		"大兴区",
		"怀柔区",
		"平谷区",
		"东城区",
		"西城区",
	}
	enc := json.NewEncoder(w)
	enc.Encode(data)
	return 0
}

func (h DistrictHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h DistrictHandler) GetHandlerMethod() string {
	return h.Method
}
