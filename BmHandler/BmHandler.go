package BmHandler

type BmHandler interface {
	GetHttpMethod() string
	GetHandlerMethod() string
}