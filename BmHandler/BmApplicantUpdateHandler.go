package BmHandler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/julienschmidt/httprouter"
)

type ApplicantUpdateHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
}

func (h ApplicantUpdateHandler) NewApplicantUpdateHandler(args ...interface{}) ApplicantUpdateHandler {
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

	return ApplicantUpdateHandler{Method: md, HttpMethod: hm, Args: ag, db: m}
}

func (h ApplicantUpdateHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h ApplicantUpdateHandler) GetHandlerMethod() string {
	return h.Method
}

func (h ApplicantUpdateHandler) UpdateApplicant(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}
	req := BmModel.Applicant{}
	json.Unmarshal(body, &req)

	err = h.db.Update(&req)

	response := map[string]interface{}{
		"status": "",
		"error":  nil,
	}

	if err == nil {
		response["status"] = "ok"
		response["error"] = err
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 0
	} else {
		response["status"] = "error"
		response["error"] = err.Error()
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}

}
