package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-Type", "text/html")
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("template parsing error: %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("templete execution error: %v", err)
		http.Error(w, "There was an error executing the template", http.StatusInternalServerError)
		return
	}

}

func handlerHome(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.gohtml")

	executeTemplate(w, tplPath)
}

func handlerContact(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")

	executeTemplate(w, tplPath)
}

func handlerFAQ(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "faq.gohtml")

	executeTemplate(w, tplPath)
}

type Router struct{}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", handlerHome)
	r.Get("/contact", handlerContact)
	r.Get("/faq", handlerFAQ)

	fmt.Println("Starting server at port 8080")

	http.ListenAndServe(":8080", r)
}
