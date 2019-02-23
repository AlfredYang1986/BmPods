package BmMiddleware

import (
	"github.com/manyminds/api2go"
	"net/http"
)

type BmMiddleware interface {
	DoMiddleware(c api2go.APIContexter, w http.ResponseWriter, r *http.Request)
}
