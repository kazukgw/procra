package procra

import (
	"net/http"
)

type Result struct {
	*http.Response
	Error error
}