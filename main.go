package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// handle the template rendering
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, req) //execute the template and write it to the response writer
}

func main() {

	var addr = flag.String("addr", ":8080", "address of the app")
	flag.Parse()
	r := newRoom() //create a new room

	http.Handle("/", &templateHandler{filename: "chat.html"}) //handle requests to the root path
	http.Handle("/room", r)                                   //handle requests to the room
	go r.run()                                                //start the room in a separate goroutine

	//start the server

	log.Println("Starting server on", *addr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
