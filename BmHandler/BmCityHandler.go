package BmHandler

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/julienschmidt/httprouter"
)

type CityHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
}

func (h CityHandler) NewCityHandler(args ...interface{}) CityHandler {
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
	return CityHandler{Method: md, HttpMethod: hm, Args: ag, db: m}
}

//TODO: load files
func (h CityHandler) AllCities(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	data := []string{
		"北京市",
	}
	enc := json.NewEncoder(w)
	enc.Encode(data)
	return 0
}

func (h CityHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h CityHandler) GetHandlerMethod() string {
	return h.Method
}
