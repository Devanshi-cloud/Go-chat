package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"time"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP implements the http.Handler interface
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t.once.Do(func() {
		templatePath := filepath.Join("templates", t.filename)
		var err error
		t.templ, err = template.ParseFiles(templatePath)
		if err != nil {
			log.Printf("Error parsing template %s: %v", templatePath, err)
			return
		}
	})

	if t.templ == nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	err := t.templ.Execute(w, req)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {

	//make every randomly generated name unique
	rand.Seed(time.Now().UnixNano())
	var addr = flag.String("addr", ":8080", "address of the app")
	flag.Parse()
	
	// Remove the old r := newRoom() and go r.run() lines - they're not needed

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //serve static files
	http.Handle("/", &templateHandler{filename: "index.html"}) //handle requests to the root path
	http.Handle("/chat", &templateHandler{filename: "chat.html"}) //handle requests to the chat page
	http.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request){
		roomName := r.URL.Query().Get("room")
		if roomName == "" {
			http.Error(w, "Room name is required", http.StatusBadRequest)
			return
		}
		realRoom := getRoom(roomName) //create a new room if it doesn't exist
		realRoom.ServeHTTP(w, r) //serve the room's HTTP handler
	})

	//start the server
	log.Println("Starting server on", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}