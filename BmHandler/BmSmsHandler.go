package BmHandler

import (
	"encoding/json"
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmSms"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

type SmsHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	d          *BmSms.BmSms
}

func (h SmsHandler) NewSmsHandler(args ...interface{}) SmsHandler {
	var d *BmSms.BmSms
	var hm string
	var md string
	var ag []string
	for i, arg := range args {
		if i == 0 {
			sts := arg.([]BmDaemons.BmDaemon)
			for _, dm := range sts {
				tp := reflect.ValueOf(dm).Interface()
				tm := reflect.ValueOf(tp).Elem().Type()
				if tm.Name() == "BmSms" {
					d = dm.(*BmSms.BmSms)
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

	return SmsHandler{Method: md, HttpMethod: hm, Args: ag, d: d}
}

type Sms struct {
	Phone  string `json:"phone" bson:"phone"`
}

func (h SmsHandler) VerifiedSmsCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	//TODO:小程序不支持patch更新，使用Function实现.
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}
	sms := Sms{}
	json.Unmarshal(body, &sms)
	err, res := h.d.SendMsg(sms.Phone)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}
	fmt.Println(res)

	return 0
}

func (h SmsHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h SmsHandler) GetHandlerMethod() string {
	return h.Method
}
