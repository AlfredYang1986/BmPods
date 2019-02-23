package BmHandler

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmPods/BmDaemons/BmMongodb"
	"github.com/alfredyang1986/BmPods/BmDaemons/BmRedis"
	"github.com/alfredyang1986/BmPods/BmModel"
	"github.com/alfredyang1986/blackmirror/bmcommon/bmsingleton"
	"github.com/alfredyang1986/blackmirror/jsonapi"
	"github.com/julienschmidt/httprouter"
)

type ApplicantHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
	rd         *BmRedis.BmRedis
}

func (h ApplicantHandler) NewApplicantHandler(args ...interface{}) ApplicantHandler {
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
	fac.RegisterModel("Applicant", &BmModel.Applicant{})
	fac.RegisterModel("BmLoginSucceed", &BmModel.BmLoginSucceed{})

	return ApplicantHandler{Method: md, HttpMethod: hm, Args: ag, db: m, rd: r}
}

func (h ApplicantHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h ApplicantHandler) GetHandlerMethod() string {
	return h.Method
}

func (h ApplicantHandler) ApplicantValidation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}
	sjson := string(body)
	rst, _ := jsonapi.FromJsonAPI(sjson)
	model := rst.(BmModel.Applicant)

	out, flag := h.checkApplicantExist(model)

	hex := md5.New()
	io.WriteString(hex, out.ID)
	token := fmt.Sprintf("%x", hex.Sum(nil))
	h.rd.PushToken(token, time.Hour*24*365)

	bmLoginSucceed := BmModel.BmLoginSucceed{
		ID:        out.ID,
		Id_:       out.Id_,
		Applicant: out,
		Token:     token,
	}

	if err != nil {
		panic(err.Error())
	}

	if flag {
		jsonapi.ToJsonAPI(&bmLoginSucceed, w)
		return 0
	} else {
		id, _ := h.registerApplicant(out)
		bmLoginSucceed.ID = id
		bmLoginSucceed.Applicant.ID = id
		jsonapi.ToJsonAPI(&bmLoginSucceed, w)
		return 0
	}
}

func (h ApplicantHandler) registerApplicant(model BmModel.Applicant) (string, error) {
	return h.db.InsertBmObject(&model)
}

func (h ApplicantHandler) checkApplicantExist(model BmModel.Applicant) (BmModel.Applicant, bool) {
	var out BmModel.Applicant
	cond := bson.M{"wechat-openid": model.WeChatOpenid}
	err := h.db.FindOneByCondition(&model, &out, cond)

	if err != nil && err.Error() == "not found" {
		return model, false
	}
	return out, true
}
