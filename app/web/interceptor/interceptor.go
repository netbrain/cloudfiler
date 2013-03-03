package interceptor

import (
	"net/http"
)

type Interceptor interface {
	Handle(http.ResponseWriter, *http.Request) bool
}
