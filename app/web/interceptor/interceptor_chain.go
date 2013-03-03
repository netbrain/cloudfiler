package interceptor

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

type InterceptorChain struct {
	interceptors []Interceptor
}

func (c *InterceptorChain) AddInterceptor(i Interceptor) {
	if c.interceptors == nil {
		c.interceptors = make([]Interceptor, 0)
	}
	c.interceptors = append(c.interceptors, i)
}

func (c *InterceptorChain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//init failure handling
	defer func() {
		if r := recover(); r != nil {
			http.Error(w, fmt.Sprintf("PANIC: %s - %s", r, debug.Stack()), http.StatusInternalServerError)
			log.Printf("PANIC: %s - %s", r, debug.Stack())
		}
	}()
	for _, i := range c.interceptors {
		if ok := i.Handle(w, r); !ok {
			break
		}
	}
}
