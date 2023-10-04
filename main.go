package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tepavcevic/go-template-server/views"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	t, err := views.Parse(filepath)
	if err != nil {
		log.Printf("template parsing error: %v", err)
		http.Error(w, "There was an error parsing the template", http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
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
