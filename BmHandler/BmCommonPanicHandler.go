package BmHandler

import (
	"fmt"
	"github.com/alfredyang1986/BmPods/BmPanic"
	"net/http"
)

type CommonPanicHandle struct {
}

func (ctm CommonPanicHandle) NewCommonPanicHandle(args ...interface{}) CommonPanicHandle {
	return CommonPanicHandle{}
}

func (ctm CommonPanicHandle) HandlePanic(rw http.ResponseWriter, r *http.Request, p interface{}) {
	fmt.Println("CommonHandlePanic接收到", p)
	BmPanic.ErrInstance().ErrorReval(p.(string), rw)
}
