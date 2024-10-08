package faust

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nokusukun/faust/docgen"
	"net/http"
)

func (api *API) Get(path string, handler func(e *Endpoint) http.HandlerFunc) *mux.Route {
	return api.Method("GET", path, handler)
}

func (api *API) Post(path string, handler func(e *Endpoint) http.HandlerFunc) *mux.Route {
	return api.Method("POST", path, handler)
}

func (api *API) Put(path string, handler func(e *Endpoint) http.HandlerFunc) *mux.Route {
	return api.Method("PUT", path, handler)
}

func (api *API) Patch(path string, handler func(e *Endpoint) http.HandlerFunc) *mux.Route {
	return api.Method("PATCH", path, handler)
}

func (api *API) Delete(path string, handler func(e *Endpoint) http.HandlerFunc) *mux.Route {
	return api.Method("DELETE", path, handler)
}

func (api *API) Method(method, path string, handler func(e *Endpoint) http.HandlerFunc) *mux.Route {
	endpoint := &Endpoint{
		EndpointInfo: EndpointInfo{
			Path:   path,
			Method: method,
		},
		Params: []IParam{},
	}
	// This is where the magic happens
	// Aka, this is where the endpoint parameter is discovered
	endpoint.httpHandler = handler(endpoint)
	api.Endpoints = append(api.Endpoints, endpoint)
	return api.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if endpoint.Use(w, r) {
			h := endpoint.httpHandler
			// we want the middleware to be executed in reverse order
			for i := len(endpoint.middlewares) - 1; i >= 0; i-- {
				h = endpoint.middlewares[i](h).ServeHTTP
			}
			h(w, r)
		}
		go endpoint.Dispose(r)
	}).Methods(method)
}

type APIContact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type APILicense struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type APIInfo struct {
	Title       string `json:"title,omitempty"`
	Summary     string `json:"summary,omitempty"`
	Description string `json:"description,omitempty"`
	//TermsOfService string     `json:"termsOfService,omitempty"`
	//Contact        APIContact `json:"contact,omitempty"`
	//License        APILicense `json:"license,omitempty"`
	Version string `json:"version,omitempty"`
}

type API struct {
	APIInfo
	isSub      bool
	Path       string      `json:"path"`
	Endpoints  []*Endpoint `json:"endpoints,omitempty"`
	Mux        *mux.Router `json:"-"`
	Subrouters []*API      `json:"subroutes,omitempty"`
	built      bool
}

func New(info ...APIInfo) *API {
	return &API{
		APIInfo: func() APIInfo {
			if len(info) > 0 {
				return info[0]
			}
			return APIInfo{}
		}(),
		Path: "/",
		Mux:  mux.NewRouter(),
	}
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !api.built {
		api.Mux.HandleFunc("/docs.json", func(w http.ResponseWriter, r *http.Request) {
			err := json.NewEncoder(w).Encode(api)
			if err != nil {
				fmt.Println(err)
			}
		}).Methods("GET")
		api.Mux.HandleFunc("/docs.html", func(w http.ResponseWriter, r *http.Request) {
			// TODO: rewrite this, this is a mess
			jsonData, _ := json.MarshalIndent(api, "", "  ")
			var apiDoc docgen.APIDoc
			err := json.Unmarshal(jsonData, &apiDoc)
			if err != nil {
				fmt.Println("Error parsing JSON:", err)
				return
			}
			html := docgen.GenerateHTML(apiDoc)

			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(200)
			w.Write([]byte(html))
		}).Methods("GET")
		api.built = true
	}
	api.Mux.ServeHTTP(w, r)
}

func (api *API) Subrouter(path string) *API {
	subApi := &API{
		Path:  path,
		isSub: true,
		Mux:   api.Mux.PathPrefix(path).Subrouter(),
	}
	api.Subrouters = append(api.Subrouters, subApi)
	return subApi
}
