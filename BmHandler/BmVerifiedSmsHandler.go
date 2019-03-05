package BmHandler

import (
	"encoding/json"
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmRedis"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmSms"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

type VerifiedSmsHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	s          *BmSms.BmSms
	r          *BmRedis.BmRedis
}

func (h VerifiedSmsHandler) NewVerifiedSmsHandler(args ...interface{}) VerifiedSmsHandler {
	var s *BmSms.BmSms
	var r *BmRedis.BmRedis
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
					s = dm.(*BmSms.BmSms)
				}
				if tm.Name() == "BmRedis" {
					r = dm.(*BmRedis.BmRedis)
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

	return VerifiedSmsHandler{Method: md, HttpMethod: hm, Args: ag, s: s, r: r}
}

func (h VerifiedSmsHandler) VerifiedSmsCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	//TODO:小程序不支持patch更新，使用Function实现.
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}
	sr := SmsRecord{}
	json.Unmarshal(body, &sr)
	err, res := h.s.VerifyCode(sr.BizId, sr.Phone)
	if err != nil {
		log.Printf("Error VerifiedSmsCode: %v", err)
		http.Error(w, "VerifiedSmsCode failed", http.StatusBadRequest)
		return 1
	}
	fmt.Println(res)
	fmt.Println(res.GetHttpContentString())
	m := make(map[string]interface{})
	err = json.Unmarshal(res.GetHttpContentBytes(), &m)
	fmt.Println(m)

	return 0
}

func (h VerifiedSmsHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h VerifiedSmsHandler) GetHandlerMethod() string {
	return h.Method
}
