package api

import (
	"fmt"
	"net/http"
)

type handlerFuncWithError func(w http.ResponseWriter, r *http.Request) error

type router struct {
	handlers map[string]handlerFuncWithError
	fallback handlerFuncWithError
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]handlerFuncWithError),
	}
}

func (r *router) get(uri string, f handlerFuncWithError) {
	r.handlers[fmt.Sprintf("%s:%s", http.MethodGet, uri)] = f
}

func (r *router) post(uri string, f handlerFuncWithError) {
	r.handlers[fmt.Sprintf("%s:%s", http.MethodPost, uri)] = f
}

func (r *router) put(uri string, f handlerFuncWithError) {
	r.handlers[fmt.Sprintf("%s:%s", http.MethodPut, uri)] = f
}

func (r *router) serveHTTP(w http.ResponseWriter, req *http.Request) error {
	if f, ok := r.handlers[fmt.Sprintf("%s:%s", req.Method, req.URL.Path)]; ok {
		return f(w, req)
	}

	return r.fallback(w, req)
}
