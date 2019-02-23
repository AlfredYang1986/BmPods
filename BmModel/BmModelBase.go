package BmModel

import "gopkg.in/mgo.v2/bson"

type BmModelBase interface {
	GetConditionsBsonM(parameters map[string][]string) bson.M
}
