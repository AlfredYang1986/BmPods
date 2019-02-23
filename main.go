package main

import (
	"fmt"
	"net/http"

	"github.com/alfredyang1986/BmPods/BmApiResolver"
	"github.com/alfredyang1986/BmPods/BmConfig"
	"github.com/alfredyang1986/BmPods/BmPodsDefine"
	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
)

func main() {
	version := "v2"
	fmt.Println("pod archi begins")
	var pod = BmPodsDefine.Pod{Name: "alfred test"}
	pod.RegisterSerFromYAML("Resources/alfredtest.yaml")

	var bmRouter BmConfig.BmRouterConfig
	bmRouter.GenerateConfig()

	addr := bmRouter.Host + ":" + bmRouter.Port
	fmt.Println("Listening on ", addr)
	api := api2go.NewAPIWithResolver(version, &BmApiResolver.RequestURL{Addr: addr})
	pod.RegisterAllResource(api)

	pod.RegisterAllFunctions(version, api)
	pod.RegisterAllMiddleware(api)

	handler := api.Handler().(*httprouter.Router)
	http.ListenAndServe(":" + bmRouter.Port, handler)

	fmt.Println("pod archi ends")
}
