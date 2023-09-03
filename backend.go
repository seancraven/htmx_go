package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	mux := http.NewServeMux()
	s := http.Server{
		Addr:         ":8002",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	mux.Handle("/", MainHandler{})
	mux.Handle("/htmx/time.html", TimeHandler{})
	mux.Handle("/htmx/ClickedGo.html", ClickGoHandler{})
	mux.Handle("/htmx/ClickedNoGo.html", NoGoHandler{})
	log.Printf("Starting server on port %v", s.Addr)
	err := s.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			panic(err)
		}
	}
}

func debugmiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%v", r)
		h.ServeHTTP(w, r)
	})
}

type TS struct {
	Time string
}

func TimeNow() TS {
	return TS{Time: time.Now().Format(time.TimeOnly)}
}

type MainHandler struct{}

func (h MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(
		"./templates/index.html",
		"./templates/htmx_time.html",
		"./templates/htmx_no_the_go.html",
	))
	t.Execute(w, TimeNow())
	//
}

type TimeHandler struct{}

func (h TimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles(
			"./templates/htmx_time.html"))
		t.Execute(w, TimeNow())
		// log.Println(TimeNow())
		//
	}
}

type TheGo struct {
	Go string
}

func TheGoload() TheGo {
	file, err := os.ReadFile("./backend.go")
	if err != nil {
		panic(err)
	}
	return TheGo{Go: string(file)}
}

type ClickGoHandler struct{}

func (h ClickGoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles(
			"./templates/htmx_the_go.html"))
		t.Execute(w, TheGoload())
	}
}

type NoGoHandler struct{}

func (h NoGoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles(
			"./templates/htmx_no_the_go.html",
		))
		t.Execute(w, nil)
	}
}
