package monocra2

import (
	"net/http"
)

type Result struct {
	*http.Response
	Error error
}
