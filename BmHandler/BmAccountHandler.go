package BmHandler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2/bson"

	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmPods/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmPods/BmDaemons/BmRedis"
	"github.com/alfredyang1986/BmPods/BmModel"

	"github.com/alfredyang1986/blackmirror/bmcommon/bmsingleton"
	"github.com/alfredyang1986/blackmirror/jsonapi"
	"github.com/alfredyang1986/blackmirror/jsonapi/jsonapiobj"
)

type AccountHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
}

func (h AccountHandler) NewAccountHandler(args ...interface{}) AccountHandler {
	var m *BmMongodb.BmMongodb
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
				if tm.Name() == "BmMongodb" {
					m = dm.(*BmMongodb.BmMongodb)
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

	//TODO: Register这边为了使用blackmirror的FromJSONAPI 2中风格迥异的Model
	fac := bmsingleton.GetFactoryInstance()
	fac.RegisterModel("Account", &BmModel.Account{})

	return AccountHandler{Method: md, HttpMethod: hm, Args: ag, db: m, rd: r}
}

func (h AccountHandler) AccountValidation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}
	sjson := string(body)
	rst, _ := jsonapi.FromJsonAPI(sjson)
	model := rst.(BmModel.Account)
	var out BmModel.Account

	cond := bson.M{"account": model.Account, "password": model.Password}

	err = h.db.FindOneByCondition(&model, &out, cond)

	jso := jsonapiobj.JsResult{}
	response := map[string]interface{}{
		"status": "",
		"result": nil,
		"error":  nil,
	}

	if err == nil && out.ID != "" {
		hex := md5.New()
		io.WriteString(hex, out.ID)
		out.Password = ""
		token := fmt.Sprintf("%x", hex.Sum(nil))
		err = h.rd.PushToken(token, time.Hour*24*365)
		out.Token = token

		response["status"] = "ok"
		response["result"] = out
		response["error"] = err

		jso.Obj = response
		enc := json.NewEncoder(w)
		enc.Encode(jso.Obj)
		return 0
	} else {
		response["status"] = "error"
		response["error"] = "账户或密码错误！"
		jso.Obj = response
		enc := json.NewEncoder(w)
		enc.Encode(jso.Obj)
		return 1
	}
}

func (h AccountHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h AccountHandler) GetHandlerMethod() string {
	return h.Method
}
