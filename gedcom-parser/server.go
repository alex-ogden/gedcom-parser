package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
)

const (
	uploadsDir = STATIC_DIR + "/uploads"
)

func runServer(srv_host string) {
	// Define endpoints
	// FileServer
	fileServer := http.FileServer(http.Dir(STATIC_DIR))
	http.Handle("/", fileServer)
	http.HandleFunc("/health", handleHealthCheck)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/render", handleRender)

	// Start the server
	log.Printf("Starting web server on: %s\n", srv_host)
	if err := http.ListenAndServe(srv_host, nil); err != nil {
		log.Fatal("Could not start server: ", err)
	}
}

// Health check
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Printf("Request method recieved was not GET\n")
		log.Printf("Request method: %s\n", r.Method)
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Got request for /health - responding with Healthy")
	fmt.Fprintf(w, "Healthy")
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Printf("Request method was not POST")
		log.Printf("Request method: %s", r.Method)
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	// Get the file from the request
	log.Println("Retrieving file from the request")
	gedcomFile, header, err := r.FormFile("gedcomFile")
	if err != nil {
		log.Fatal(err)
	}

	// Read file data into memory
	log.Println("Reading file data into memory")
	fileData, err := io.ReadAll(gedcomFile)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Parsing GEDCOM data")
	returnPeople, err := parseFile(fileData, header)
	if err != nil {
		log.Fatal(err)
	}

	encodedPeopleList := base64.StdEncoding.EncodeToString([]byte(returnPeople))
	redirectUrl := fmt.Sprintf("/render?peopleList=%s", encodedPeopleList)
	log.Println("Sending request to render page")
	http.Redirect(w, r, redirectUrl, 301)
}

func handleRender(w http.ResponseWriter, r *http.Request) {
	// Get base64 encoded list from request
	parameters := r.URL.Query()["peopleList"]
	encodedPeopleList := parameters[0]
	peopleList, err := base64.StdEncoding.DecodeString(encodedPeopleList)
	if err != nil {
		log.Fatal(err)
	}

	tmpl := template.Must(template.ParseFiles(STATIC_DIR + "/render.html"))
	data := RenderData{
		PeopleList: string(peopleList),
	}

	log.Printf("Rendering results page\n")
	if err := tmpl.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}
