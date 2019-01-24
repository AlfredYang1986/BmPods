package BmMiddleware

import (
	"github.com/alfredyang1986/BmPods/BmDaemons"
	"github.com/alfredyang1986/BmPods/BmDaemons/BmRedis"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
)

type CheckTokenMiddleware struct {
	Args        []string
	rd          *BmRedis.BmRedis
}

func (ctm CheckTokenMiddleware) NewCheckTokenMiddleware(args ...interface{}) CheckTokenMiddleware {
	var r *BmRedis.BmRedis
	var ag []string
	for i, arg := range args {
		if i == 0 {
			sts := arg.([]BmDaemons.BmDaemon)
			for _, dm := range sts {
				tp := reflect.ValueOf(dm).Interface()
				tm := reflect.ValueOf(tp).Elem().Type()
				if tm.Name() == "BmRedis" {
					r = dm.(*BmRedis.BmRedis)
				}
			}
		} else if i == 1 {
			lst := arg.([]string)
			for _, str := range lst {
				ag = append(ag, str)
			}
		} else {
		}
	}

	return CheckTokenMiddleware{Args: ag, rd: r}
}

func (ctm CheckTokenMiddleware) DoMiddleware(c api2go.APIContexter, w http.ResponseWriter, r *http.Request) {
	//path := r.URL.Path
	//if strings.HasSuffix(path, "AccountValidation") || strings.HasSuffix(path, "ApplicantValidation") {
	//	fmt.Println("login from ", path)
	//} else {
	//	fmt.Println("Check Token!")
	//	auth := r.Header.Get("Authorization")
	//	arr := strings.Split(auth, " ")
	//	if len(arr) < 2 || arr[0] != "bearer" {
	//		panic("Auth Failed!")
	//	}
	//	token := arr[1]
	//	err := ctm.rd.CheckToken(token)
	//	if err != nil {
	//		panic(err.Error())
	//	}
	//}
}
