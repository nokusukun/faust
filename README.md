# Faust
A net/http compatible router and middleware for that has FastAPI-like features.

*This is still in prototyping stage*


```go
package main

import (
	"encoding/json"
	"github.com/nokusukun/faust"
	"github.com/nokusukun/faust/param"
	"net/http"
)

/*
from typing import Union
from fastapi import FastAPI

app = FastAPI()

@app.get("/")
def read_root():
	return {"Hello": "World"}

@app.get("/items/{item_id}")
def read_item(item_id: int, q: Union[str, None] = None):
	return {"item_id": item_id, "q": q}
*/

func main() {
	app := faust.New()
	
	app.Get("/", func(e *faust.Endpoint) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]any{
				"Hello": "World",
			})
		}
	})

	app.Get("/items/{item_id}", func(e *faust.Endpoint) http.HandlerFunc {
		itemId := param.Path[string](e, "item_id", param.Info{
			Description: "ID of the item, must start with item-",
		})

		q := param.Query[string](e, "q", param.Info{
			Description: "Query string",
		}).Optional()

		return func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]any{
				"item_id": itemId.Value(r),
				"q":       q.Value(r),
			})
		}
	})

	err := http.ListenAndServe(":8081", app)
	if err != nil {
		panic(err)
	}
}


```