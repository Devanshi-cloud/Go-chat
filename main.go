package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"strings"
)

type templateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
}

type roomHandler struct {
	rooms map[string]*room
	mutex sync.RWMutex
}

//handling the template from our server

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("template", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func (rh *roomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract room name from query parameter
	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		http.Error(w, "Room name required", http.StatusBadRequest)
		return
	}

	// Get or create room
	rh.mutex.Lock()
	room, exists := rh.rooms[roomName]
	if !exists {
		room = newRoom(roomName)
		rh.rooms[roomName] = room
		go room.run()
	}
	rh.mutex.Unlock()

	// Handle WebSocket upgrade
	room.ServeHTTP(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "application address")
	flag.Parse()

	// Initialize room handler
	rh := &roomHandler{
		rooms: make(map[string]*room),
	}

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Handle routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// Serve index page
			tmpl := template.Must(template.ParseFiles(filepath.Join("template", "index.html")))
			tmpl.Execute(w, nil)
		} else if strings.HasPrefix(r.URL.Path, "/chat") {
			// Serve chat page
			tmpl := template.Must(template.ParseFiles(filepath.Join("template", "chat.html")))
			tmpl.Execute(w, nil)
		} else {
			http.NotFound(w, r)
		}
	})

	http.Handle("/room", rh)

	log.Println("Starting web server on", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Listen And Serve:", err)
	}
}