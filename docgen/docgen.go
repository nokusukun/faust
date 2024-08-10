package docgen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
)

type Schema struct {
	Type   string `json:"type"`
	Format string `json:"format"`
}

type Parameter struct {
	In          string `json:"in"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      Schema `json:"schema"`
}

type Endpoint struct {
	Method      string      `json:"method"`
	Path        string      `json:"path"`
	Description string      `json:"description"`
	Parameters  []Parameter `json:"parameters"`
}

type Subroute struct {
	Path      string     `json:"path"`
	Endpoints []Endpoint `json:"endpoints"`
}

type APIDoc struct {
	Title     string     `json:"title"`
	Summary   string     `json:"summary"`
	Version   string     `json:"version"`
	Path      string     `json:"path"`
	Subroutes []Subroute `json:"subroutes"`
	Endpoints []Endpoint `json:"endpoints"`
}

func GenerateHTML(apiDoc APIDoc) string {
	if apiDoc.Title == "" {
		apiDoc.Title = " Faust API"
	}
	if apiDoc.Summary == "" {
		apiDoc.Summary = "API Documentation for the Faust API"
	}
	if apiDoc.Version == "" {
		apiDoc.Version = "1.0.0"
	}
	t := template.Must(template.New("apiDoc").Parse(tmpl))
	var result bytes.Buffer
	err := template.Must(t.Clone()).Execute(&result, apiDoc)
	if err != nil {
		panic(err)
	}
	return result.String()
}

func main() {
	jsonData := `{
		"title": "Users API",
		"summary": "Sample API for managing users",
		"version": "0.0.1",
		"path": "/",
		"subroutes": [
			{
				"path": "/users",
				"endpoints": [
					{
						"method": "GET",
						"path": "/{id}",
						"description": "Fetch a user by ID",
						"parameters": [
							{
								"in": "path",
								"name": "id",
								"description": "The ID of the user to fetch",
								"schema": {
									"type": "uint",
									"format": "uint"
								}
							}
						]
					},
					{
						"method": "POST",
						"path": "/",
						"description": "Create a new user",
						"parameters": [
							{
								"in": "jsonbody",
								"name": "payload",
								"description": "The payload to create a new user",
								"schema": {
									"type": "struct",
									"format": "struct"
								}
							}
						]
					},
					{
						"method": "PUT",
						"path": "/{id}",
						"description": "Update a user by ID",
						"parameters": [
							{
								"in": "path",
								"name": "id",
								"description": "The ID of the user to update",
								"schema": {
									"type": "uint",
									"format": "uint"
								}
							},
							{
								"in": "jsonbody",
								"name": "payload",
								"description": "The payload to update a user",
								"schema": {
									"type": "struct",
									"format": "struct"
								}
							}
						]
					},
					{
						"method": "DELETE",
						"path": "/{id}",
						"description": "Delete a user by ID",
						"parameters": [
							{
								"in": "path",
								"name": "id",
								"description": "The ID of the user to delete",
								"schema": {
									"type": "uint",
									"format": "uint"
								}
							}
						]
					}
				]
			}
		]
	}`

	var apiDoc APIDoc
	err := json.Unmarshal([]byte(jsonData), &apiDoc)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	html := GenerateHTML(apiDoc)
	f, err := os.Create("api_doc.html")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(html)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("API documentation generated successfully!")
}
