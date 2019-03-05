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
	"math/rand"
	"net/http"
	"reflect"
	"time"
)

type GenerateSmsHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	s          *BmSms.BmSms
	r          *BmRedis.BmRedis
}

func (h GenerateSmsHandler) NewGenerateSmsHandler(args ...interface{}) GenerateSmsHandler {
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

	return GenerateSmsHandler{Method: md, HttpMethod: hm, Args: ag, s: s, r: r}
}

type Sms struct {
	Phone string `json:"phone" bson:"phone"`
}

type SmsRecord struct {
	BizId string `json:"biz-id" bson:"biz-id"`
	Phone string `json:"phone" bson:"phone"`
	Code  string `json:"code" bson:"code"`
}

func (h GenerateSmsHandler) GenerateSmsCode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	//TODO:小程序不支持patch更新，使用Function实现.
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}
	sms := Sms{}
	response := map[string]interface{}{
		"status": "",
		"result": nil,
		"error":  nil,
	}
	json.Unmarshal(body, &sms)
	err, res := h.s.SendMsg(sms.Phone, GenerateRandNumber())
	if err != nil {
		log.Printf("Error SendMsg: %v", err)
		response["status"] = "error"
		response["error"] = err
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(res.GetHttpContentBytes(), &m)
	fmt.Println(m)


	bizId, ok := m["BizId"]
	if ok {
		sr := SmsRecord{}
		sr.Phone = sms.Phone
		sr.BizId = bizId.(string)

		sr.Code = "ok"
		response["status"] = "ok"
		response["result"] = sr
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 0
	} else {
		response["status"] = "error"
		response["error"] = "no BizId found! 同一手机号频繁调用!"
		enc := json.NewEncoder(w)
		enc.Encode(response)
		return 1
	}

}

func (h GenerateSmsHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h GenerateSmsHandler) GetHandlerMethod() string {
	return h.Method
}

func GenerateRandNumber() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rst := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	return rst
}
