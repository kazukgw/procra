package procra

import (
	"net/http"
)

type ProxyState interface {
	String() string
	Fetch(*TargetURL) (*http.Response, error)
	HandleResult(*Proxy, *Result) error
}
