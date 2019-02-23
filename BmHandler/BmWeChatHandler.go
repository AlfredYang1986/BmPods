package BmHandler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/alfredyang1986/BmPods/BmDaemons"
	"github.com/alfredyang1986/BmPods/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/alfredyang1986/blackmirror/jsonapi/jsonapiobj"
	"github.com/julienschmidt/httprouter"
)

type WeChatHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
}

func (h WeChatHandler) NewWeChatHandler(args ...interface{}) WeChatHandler {
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
	return WeChatHandler{Method: md, HttpMethod: hm, Args: ag, db: m}
}

func (h WeChatHandler) GetWeChatInfo(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")
	var input BmModel.WeChatInfo
	json.NewDecoder(r.Body).Decode(&input)

	jso := jsonapiobj.JsResult{}
	response := map[string]interface{}{
		"status": "",
		"result": nil,
		"error":  nil,
	}

	if input.ServerName == "" {
		ems := "no server-name found"
		response["status"] = "error"
		response["error"] = ems
		jso.Obj = response
		enc := json.NewEncoder(w)
		enc.Encode(jso.Obj)
		return 0
	}
	if input.Code == "" {
		ems := "no code found"
		response["status"] = "error"
		response["error"] = ems
		jso.Obj = response
		enc := json.NewEncoder(w)
		enc.Encode(jso.Obj)
		return 0
	}

	switch input.ServerName {
	case "dongda":
		input.AppId = "wx6129e48a548c52b8"
		input.Secret = "b250e875e51a931e2ae3a49ff450bc3c"
	case "pacee":
		input.AppId = "wx79138b2ee5288cc2"
		input.Secret = "c2637375412cfa97c9e127b4cde30c5c"
	}

	originUrl := "https://api.weixin.qq.com/sns/jscode2session?appid="
	url := strings.Join([]string{originUrl, input.AppId, "&secret=", input.Secret, "&js_code=", input.Code, "&grant_type=authorization_code"}, "")
	resp, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	m := make(map[string]interface{})
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		panic(err.Error())
	}

	if v, ok := m["openid"]; ok {
		input.OpenId = v.(string)
	} else {
		response["status"] = "error"
		response["error"] = m
		jso.Obj = response
		enc := json.NewEncoder(w)
		enc.Encode(jso.Obj)
		return 0
	}
	if v, ok := m["session_key"]; ok {
		input.SessionKey = v.(string)
	} else {
		response["status"] = "error"
		response["error"] = m
		jso.Obj = response
		enc := json.NewEncoder(w)
		enc.Encode(jso.Obj)
		return 0
	}

	response["status"] = "ok"
	response["result"] = input
	response["error"] = nil
	jso.Obj = response
	enc := json.NewEncoder(w)
	enc.Encode(jso.Obj)
	return 1
}

func (h WeChatHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h WeChatHandler) GetHandlerMethod() string {
	return h.Method
}
