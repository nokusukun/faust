package param

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nokusukun/faust"
	cmap "github.com/orcaman/concurrent-map"
	"net/http"
	"reflect"
	"strconv"
)

func reqToId(r *http.Request) string {
	return fmt.Sprintf("%p", r)
}

func Query[T any](e *faust.Endpoint, name string, paramInfo ...Info) *EndpointParam[T] {
	return Param[T]("query", e, name, paramInfo...)
}

func Path[T any](e *faust.Endpoint, name string, paramInfo ...Info) *EndpointParam[T] {
	return Param[T]("path", e, name, paramInfo...)
}

func Body[T any](e *faust.Endpoint, name string, paramInfo ...Info) *EndpointParam[T] {
	return Param[T]("body", e, name, paramInfo...)
}

func Json[T any](e *faust.Endpoint, name string, paramInfo ...Info) *EndpointParam[T] {
	return Param[T]("jsonbody", e, name, paramInfo...)
}

func Header[T any](e *faust.Endpoint, name string, paramInfo ...Info) *EndpointParam[T] {
	return Param[T]("header", e, name, paramInfo...)
}

func Form[T any](e *faust.Endpoint, name string, paramInfo ...Info) *EndpointParam[T] {
	return Param[T]("form", e, name, paramInfo...)
}

func Param[T any](ptype string, e *faust.Endpoint, name string, paramInfo ...Info) *EndpointParam[T] {
	tType := reflect.TypeOf(new(T)).Elem()
	switch tType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.String, reflect.Struct:
		break
	default:
		panic("unsupported type")
	}

	param := &EndpointParam[T]{
		outType: tType,
		values:  cmap.New(),
	}
	if len(paramInfo) > 0 {
		param.parameterInfo.Info = paramInfo[0]
	}
	param.parameterInfo.In = ptype
	param.parameterInfo.Name = name
	param.Schema.Type = tType.Kind().String()
	param.Schema.Format = tType.Kind().String()
	e.Params = append(e.Params, param)
	return param
}

type ParameterSchema struct {
	Type   string `json:"type,omitempty"`
	Format string `json:"format,omitempty"`
}

type Info struct {
	Description string `json:"description,omitempty"`
	Optional    bool   `json:"optional,omitempty"`
}

type parameterInfo struct {
	In   string `json:"in,omitempty"`
	Name string `json:"name,omitempty"`
	Info
	Schema ParameterSchema `json:"schema"`
}

type EndpointParam[T any] struct {
	parameterInfo
	outType   reflect.Type
	values    cmap.ConcurrentMap
	validator []func(T) error
}

func (e *EndpointParam[T]) Dispose(r *http.Request) {
	reqId := reqToId(r)
	e.values.Remove(reqId)
}

func (e *EndpointParam[T]) Validate(validateFunc ...func(T) error) *EndpointParam[T] {
	e.validator = validateFunc
	return e
}

func (e *EndpointParam[T]) Description(desc string) *EndpointParam[T] {
	e.parameterInfo.Description = desc
	return e
}

func (e *EndpointParam[T]) Optional() *EndpointParam[T] {
	e.parameterInfo.Optional = true
	return e
}

func (e *EndpointParam[T]) Use(r *http.Request) error {
	var err error
	val, err := e.ValueWithError(r)
	if err != nil {
		return err
	}
	if e.validator != nil {
		for _, validate := range e.validator {
			if err := validate(val); err != nil {
				return fmt.Errorf("%s (%v:%v)", err.Error(), e.In, e.Name)
			}
		}
	}
	e.values.Set(reqToId(r), val)
	return nil
}

func (e *EndpointParam[T]) ValueWithError(r *http.Request) (T, error) {
	if _, ok := e.values.Get(reqToId(r)); ok {
		return e.Value(r), nil
	}

	var t T
	var value string
	switch e.parameterInfo.In {
	case "query":
		if !e.Info.Optional && !r.URL.Query().Has(e.parameterInfo.Name) {
			return t, fmt.Errorf("missing required parameter %s", e.parameterInfo.Name)
		}
		value = r.URL.Query().Get(e.parameterInfo.Name)
	case "path":
		var exists bool
		value, exists = mux.Vars(r)[e.parameterInfo.Name]
		if !e.Info.Optional && !exists {
			return t, fmt.Errorf("missing required parameter %s", e.parameterInfo.Name)
		}
	case "header":
		v, exists := r.Header[e.parameterInfo.Name]
		if !e.Info.Optional && !exists {
			return t, fmt.Errorf("missing required parameter %s", e.parameterInfo.Name)
		}
		value = v[0]
	case "form":
		if !e.Info.Optional && !r.Form.Has(e.parameterInfo.Name) {
			return t, fmt.Errorf("missing required parameter %s", e.parameterInfo.Name)
		}
		value = r.FormValue(e.parameterInfo.Name)
	case "body":
		var body []byte
		c, err := r.Body.Read(body)
		if err != nil {
			return t, err
		}
		if c == 0 {
			return t, fmt.Errorf("missing required body %s", e.parameterInfo.Name)
		}
		value = string(body)
	case "jsonbody":
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			return t, err
		}
		return t, nil
	default:
		return t, nil
	}

	switch e.outType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parseInt, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return t, err
		}
		return reflect.ValueOf(parseInt).Convert(reflect.TypeOf(t)).Interface().(T), nil
	case reflect.Float32, reflect.Float64:
		parseFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return t, err
		}
		return reflect.ValueOf(parseFloat).Interface().(T), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parseUint, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return t, err
		}
		return reflect.ValueOf(parseUint).Convert(reflect.TypeOf(t)).Interface().(T), nil
	case reflect.String:
		return reflect.ValueOf(value).Interface().(T), nil
	}
	return t, nil
}

func (e *EndpointParam[T]) Value(r *http.Request) T {
	//return e.values[r]
	val, _ := e.values.Get(fmt.Sprintf("%p", r))
	return val.(T)
}
