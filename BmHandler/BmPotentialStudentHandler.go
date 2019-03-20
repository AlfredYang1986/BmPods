package BmHandler

import (
	"encoding/json"
	"github.com/alfredyang1986/BmPods/BmModel"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/julienschmidt/httprouter"
)

type PotentialStudentHandler struct {
	Method     string
	HttpMethod string
	Args       []string
	db         *BmMongodb.BmMongodb
}

func (h PotentialStudentHandler) NewPotentialStudentHandler(args ...interface{}) PotentialStudentHandler {
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

	return PotentialStudentHandler{Method: md, HttpMethod: hm, Args: ag, db: m}
}

func (h PotentialStudentHandler) GetHttpMethod() string {
	return h.HttpMethod
}

func (h PotentialStudentHandler) GetHandlerMethod() string {
	return h.Method
}

type PotentialStudent struct {
	BrandId      string  `json:"brand-id" bson:"brand-id"`
	Name         string  `json:"name" bson:"name"`
	School       string  `json:"school" bson:"school"`
	Gender       float64 `json:"gender" bson:"gender"`
	Grade        string  `json:"grade" bson:"grade"`
	Dob          float64 `json:"dob" bson:"dob"`
	GuardianName string  `json:"guardian-name" bson:"guardian-name"`
	Contact      string  `json:"contact" bson:"contact"`
	RelationShip string  `json:"relation-ship" bson:"relation-ship"`
}

func (h PotentialStudentHandler) AddPotentialStudent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) int {
	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return 1
	}
	potentialStudent := PotentialStudent{}
	json.Unmarshal(body, &potentialStudent)

	stud := BmModel.Student{
		BrandId: potentialStudent.BrandId,
		Name:    potentialStudent.Name,
		School:  potentialStudent.School,
		Gender:  potentialStudent.Gender,
		Grade:   potentialStudent.Grade,
		Dob:     potentialStudent.Dob,
	}

	applicant, err := h.checkApplicantExist(potentialStudent)
	if err != nil {
		panic(err.Error())
	}
	kid := BmModel.Kid{
		Name:         potentialStudent.Name,
		Gender:       potentialStudent.Gender,
		Dob:          potentialStudent.Dob,
		GuardianRole: potentialStudent.RelationShip,
		ApplicantID:  applicant.ID,
	}
	kidID, err := h.db.InsertBmObject(&kid)
	if err != nil {
		panic(err.Error())
	}
	stud.KidID = kidID

	guardian, err := h.checkGuardianExist(potentialStudent)
	if err != nil {
		panic(err.Error())
	}
	stud.GuardiansIDs = []string{guardian.ID}
	stud.Contact = guardian.Contact
	stud.RegDate = float64(time.Now().UnixNano() / 1e6)
	stud.CreateTime = stud.RegDate

	studID, err := h.db.InsertBmObject(&stud)
	if err != nil {
		panic(err.Error())
	}
	stud.ID = studID

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

func (h PotentialStudentHandler) checkGuardianExist(ps PotentialStudent) (BmModel.Guardian, error) {
	var out BmModel.Guardian
	cond := bson.M{"contact": ps.Contact}
	err := h.db.FindOneByCondition(&out, &out, cond)

	if err != nil && err.Error() == "not found" {
		out.Name = ps.GuardianName
		out.Contact = ps.Contact
		out.RelationShip = ps.RelationShip
		out.BrandId = ps.BrandId
		out.RegDate = float64(time.Now().UnixNano() / 1e6)
		id, err := h.db.InsertBmObject(&out)
		if err != nil {
			panic(err.Error())
		}
		out.ID = id
		return out, err
	}
	return out, err
}

func (h PotentialStudentHandler) checkApplicantExist(ps PotentialStudent) (BmModel.Applicant, error) {
	var out BmModel.Applicant
	cond := bson.M{"regi-phone": ps.Contact}
	err := h.db.FindOneByCondition(&out, &out, cond)

	if err != nil && err.Error() == "not found" {
		out.Name = ps.GuardianName
		out.RegisterPhone = ps.Contact
		out.WeChatBindPhone = ps.Contact
		id, err := h.db.InsertBmObject(&out)
		if err != nil {
			panic(err.Error())
		}
		out.ID = id
		return out, err
	}
	return out, err
}
