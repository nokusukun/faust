# Faust


`faust` is a net/http compatible router and middleware for that has FastAPI-like features. It provides a convenient way to define HTTP endpoints, handle requests, and validate parameters, all while offering automatic API documentation generation.


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

## Features

- **Simple API Definition**: Easily define RESTful API endpoints for various HTTP methods (GET, POST, PUT, PATCH, DELETE).
- **Parameter Handling**: Support for handling query parameters, path parameters, headers, form data, and JSON bodies with type safety and validation.
- **Middleware Support**: Chain middlewares for each endpoint to handle cross-cutting concerns like authentication, logging, etc.
- **Automatic Documentation**: Automatically generate API documentation in both JSON and HTML formats.
- **Subrouters**: Organize your API into subrouters to create modular and structured routes.

## Installation

To install `faust`, run:

```bash
go get github.com/nokusukun/faust
```

## Usage

### Initializing the API

```go
package main

import (
    "net/http"
    "github.com/nokusukun/faust"
)

func main() {
    api := faust.New(faust.APIInfo{
        Title:       "My API",
        Version:     "1.0.0",
        Description: "This is my API built using Faust.",
    })

    api.Get("/hello", func(e *faust.Endpoint) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            w.Write([]byte("Hello, World!"))
        }
    })

    http.ListenAndServe(":8080", api)
}
```

### Defining Endpoints

Faust provides convenient methods to define endpoints for different HTTP methods:

```go
api.Get("/path", handlerFunc)
api.Post("/path", handlerFunc)
api.Put("/path", handlerFunc)
api.Patch("/path", handlerFunc)
api.Delete("/path", handlerFunc)
```

Each handler function is provided with an `Endpoint` object to manage parameters, middlewares, and more.

### Handling Parameters

Faust provides utilities to easily handle and validate parameters from various sources:

```go
import "github.com/nokusukun/faust/param"

func ExampleHandler(e *faust.Endpoint) http.HandlerFunc {
    queryParam := param.Query[string](e, "example_query").Optional().Description("An example query parameter")

    return func(w http.ResponseWriter, r *http.Request) {
        val := queryParam.Value(r)
        w.Write([]byte("Query parameter value: " + val))
    }
}
```

### Middlewares

You can add middlewares to your endpoints:

```go
var AuthMiddleware mux.MiddlewareFunc = func(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("AuthMiddleware")
        auth := r.Header.Get("Authorization")
        if auth == "" {
            ReturnJSON(w, map[string]any{
                "error": "Authorization header is required",
            }, 401)
            return
        }
        if !strings.HasPrefix(auth, "Basic ") {
            ReturnJSON(w, map[string]any{
                "error": "Authorization header must be a Basic auth",
			}, 401)
            return
        }
        if auth != fmt.Sprintf("Basic %v:%v", config.USERNAME, config.PASSWORD) {
            ReturnJSON(w, map[string]any{
                "error": "Invalid credentials",
            }, 401)
            return
        }
        handler.ServeHTTP(w, r)
    })
}

e.Middlewares(LoggingMiddleware, AuthMiddleware)
```

### Subrouters

To organize your routes, you can use subrouters:

```go
subApi := api.Subrouter("/sub")
subApi.Get("/example", ExampleHandler)
```

### Generating Documentation

Faust can automatically generate API documentation in JSON and HTML formats:

- **JSON Documentation**: Accessible at `/docs.json`
- **HTML Documentation**: Accessible at `/docs.html`

## Example

Here’s a more comprehensive example:

```go
package main

import (
    "net/http"
    "github.com/nokusukun/faust"
    "github.com/nokusukun/faust/param"
)

func main() {
    api := faust.New(faust.APIInfo{
        Title:       "My API",
        Version:     "1.0.0",
        Description: "This is my API built using Faust.",
    })

    api.Get("/greet", func(e *faust.Endpoint) http.HandlerFunc {
        nameParam := param.Query[string](e, "name").Optional().Description("The name of the person to greet")

        return func(w http.ResponseWriter, r *http.Request) {
            name := nameParam.Value(r)
            if name == "" {
                name = "World"
            }
            w.Write([]byte("Hello, " + name + "!"))
        }
    })

    http.ListenAndServe(":8080", api)
}
```

## License

~~This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.~~

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.
