package BmHandler

import (
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmPods/BmDaemons/BmMongodb"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"reflect"
)

type UpdateHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
}

func (h UpdateHandler) NewUpdateHandler(args ...interface{}) UpdateHandler {
	var m *BmMongodb.BmMongodb
	var hm string
	var md string
	var ag []string
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

	return UpdateHandler{Method: md, HttpMethod: hm, Args: ag, db: m}
}

func (h UpdateHandler) UpdateFunction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	//TODO:小程序不支持patch更新，使用Function实现.

	return 0
}

func (h UpdateHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h UpdateHandler) GetHandlerMethod() string {
	return h.Method
}
