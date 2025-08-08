package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {

	once sync.Once
	filename string
	templ *template.Template
}

//handling the template from our server

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("template", t.filename)))
	})
	t.templ.Execute(w, nil)
}

func main() {
	var addr = flag.String("addr", ":8080", "application address")
	flag.Parse()
	r := newRoom()

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.Handle("/", &templateHandler{filename: "chat.html"})

	http.Handle("/room", r)

	go r.run()

	//start the room

	log.Println("Starting web server on", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Listen And Serve:", err)
	}
}