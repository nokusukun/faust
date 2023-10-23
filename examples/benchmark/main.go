package main

import (
	main2 "api-project-a"
	"api-project-a/param"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
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

type Item struct {
	ItemId  int64  `json:"item_id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func main() {
	app := main2.New()
	app.Get("/", func(e *main2.Endpoint) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"Hello": "World"}`))
		}
	})

	app.Get("/items/{item_id}", func(e *main2.Endpoint) http.HandlerFunc {
		itemId := param.Path[string](e, "item_id")
		q := param.Query[string](e, "q")
		return func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"itemId": ` + itemId.Value(r) + `, "q": ` + q.Value(r) + `}`))
		}
	})

	lock := sync.RWMutex{}
	i := map[string]Item{}
	app.Post("/items/{item_id}", func(e *main2.Endpoint) http.HandlerFunc {
		itemId := param.Path[string](e, "item_id")
		item := param.Json[Item](e, "item").Required()
		return func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
			lock.Lock()
			i[itemId.Value(r)] = item.Value(r)
			lock.Unlock()
			w.Write([]byte(`OK`))
		}
	})

	app.Get("/items/{item_id}/content", func(e *main2.Endpoint) http.HandlerFunc {
		itemId := param.Path[string](e, "item_id")
		return func(w http.ResponseWriter, r *http.Request) {
			lock.RLock()
			json.NewEncoder(w).Encode(i[itemId.Value(r)])
			lock.RUnlock()
		}
	})

	go Benchmark()

	err := http.ListenAndServe(":8081", app)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Define a constant string of characters to choose from
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Create a new random source using the current time as a seed
var src = rand.NewSource(time.Now().UnixNano())

// Create a new random generator using the source
var rnd = rand.New(src)

// Define a function that takes an integer n as an argument
// and returns a random string of length n
func randomString(n int) string {
	// Create a byte slice of length n
	b := make([]byte, n)
	// Loop over the byte slice
	for i := range b {
		// For each index, assign a random byte from the charset
		b[i] = charset[rnd.Intn(len(charset))]
	}
	// Convert the byte slice to a string and return it
	return string(b)
}

var client = http.Client{}

func Scream(i int) {
	fmt.Println("Scream", i)
	item := Item{
		ItemId:  int64(rand.Int()),
		Name:    randomString(rand.Intn(100)),
		Content: randomString(rand.Intn(100)),
	}
	var payload io.ReadWriter
	payload = &bytes.Buffer{}
	err := json.NewEncoder(payload).Encode(item)
	if err != nil {
		fmt.Printf("error encoding payload: %v\n", err)
	}
	req, err := http.NewRequest("POST", "http://localhost:8081/items/"+fmt.Sprintf("%d", item.ItemId), payload)
	if err != nil {
		fmt.Printf("error creating request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error sending request: %v\n", err)
	}
	resp.Body.Close()

	req, err = http.NewRequest("GET", "http://localhost:8081/items/"+fmt.Sprintf("%d", item.ItemId)+"/content", nil)
	if err != nil {
		fmt.Printf("error creating request: %v\n", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("error sending request: %v\n", err)
	}
	reqItem := Item{}
	err = json.NewDecoder(resp.Body).Decode(&reqItem)
	if err != nil {
		fmt.Printf("error decoding response: %v\n", err)
	}
	resp.Body.Close()
	if reqItem.ItemId != item.ItemId {
		fmt.Printf("item id mismatch: %d != %d\n", reqItem.ItemId, item.ItemId)
	}
	if reqItem.Name != item.Name {
		fmt.Printf("item name mismatch: %s != %s\n", reqItem.Name, item.Name)
	}
	if reqItem.Content != item.Content {
		fmt.Printf("item content mismatch: %s != %s\n", reqItem.Content, item.Content)
	}
}

func Benchmark() {
	fmt.Println("Running benchmark in 1 sec")
	time.Sleep(time.Second)
	numWorkers := 2
	for i := 0; i < numWorkers; i++ {
		go func() {
			for i := 0; ; i++ {
				now := time.Now()
				Scream(i)
				fmt.Println(time.Since(now))
			}
		}()
	}
}
