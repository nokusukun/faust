package main

import (
	"fmt"
	"github.com/nokusukun/faust"
	"github.com/nokusukun/faust/param"
	"net/http"
)

func ShouldNotBeAdmin(s string) error {
	if s == "admin" {
		return fmt.Errorf("admin is not allowed")
	}
	return nil
}

func main() {

	myApi := faust.New(faust.APIInfo{
		Title:       "Sample API",
		Summary:     "Sample API Summary",
		Description: "Description of the API",
		Version:     "v0.0.1",
	})

	myApi.Get("/api/status", func(e *faust.Endpoint) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		}
	})

	users := myApi.Subrouter("/api/v1/users")

	users.Get("/{id}", func(e *faust.Endpoint) http.HandlerFunc {
		e.Name("Get User Info")
		e.Description("Get user info")
		id := param.Path[int64](e, "id", param.Info{
			Description: "ID of the user",
		})

		return func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fmt.Sprintf("ID: %v", id.Value(r))))
		}
	})

	userOperations := users.Subrouter("/operations")

	userOperations.Get("/delete", func(e *faust.Endpoint) http.HandlerFunc {
		e.Name("delete")
		e.Description("Delete user")
		_ = param.Query[string](e, "name", param.Info{
			Description: "name of the user",
		}).Validate(ShouldNotBeAdmin)
		_ = param.Query[int64](e, "timestamp", param.Info{
			Description: "timestamp of the operation",
		})

		return func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		}
	})

	fmt.Println("Listening on :8081")
	err := http.ListenAndServe(":8081", myApi)
	if err != nil {
		fmt.Println(err)
		return
	}

	select {}
}
