package faust

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type EndpointInfo struct {
	Method      string `json:"method,omitempty"`
	Path        string `json:"path,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type Endpoint struct {
	EndpointInfo
	Params      []IParam `json:"parameters,omitempty"`
	middlewares []mux.MiddlewareFunc
	httpHandler http.HandlerFunc
	OnError     func(w http.ResponseWriter, r *http.Request, err error) `json:"-"`
}

func (e *Endpoint) Middlewares(middlewares ...mux.MiddlewareFunc) *Endpoint {
	e.middlewares = append(e.middlewares, middlewares...)
	return e
}

func (e *Endpoint) Name(name string) *Endpoint {
	e.EndpointInfo.Name = name
	return e
}

func (e *Endpoint) Description(name string) *Endpoint {
	e.EndpointInfo.Description = name
	return e
}

func (e *Endpoint) UseErr(r *http.Request) error {
	for _, param := range e.Params {
		if err := param.Use(r); err != nil {
			return err
		}
	}
	return nil
}

func (e *Endpoint) Use(w http.ResponseWriter, r *http.Request) bool {
	err := e.UseErr(r)
	if err != nil {
		if e.OnError != nil {
			e.OnError(w, r, err)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(map[string]any{
				"error": err.Error(),
				"type":  "validation_error",
			})
		}
		return false
	}

	return true
}

func (e *Endpoint) Dispose(r *http.Request) {
	for _, param := range e.Params {
		param.Dispose(r)
	}
}

type IParam interface {
	Use(r *http.Request) error
	Dispose(r *http.Request)
}
