package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func handlerHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tpl, err := template.ParseFiles("templates/home.gohtml")
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

func handlerContact(w http.ResponseWriter, r *http.Request) {
	contactID := chi.URLParam(r, "contactID")
	fmt.Println(contactID)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To contact me send an email to <a href=\"mailto:djordje.tepa@gmail.com\">djordje.tepa@gmail.com</a>.</p>")
}

func handlerFAQ(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(
		w,
		`<h1>FAQ</h1>
		<p>Q: Is there a free version?</p>
		<p>A: Yes! We offer a free trial for 30 days on any period plans.</p>
		<br>
		<p>Q: What are Your support hours?</p>
		<p>A: Yes! We offer a free trial for 30 days on any period plans.</p>
		<br>
		<p>Q: Is there a free version?</p>
		<p>A: Yes! We offer a free trial for 30 days on any period plans.</p>
		`,
	)
}

type Router struct{}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", handlerHome)
	r.Get("/contact/{contactID}", handlerContact)
	r.Get("/faq", handlerFAQ)

	fmt.Println("Starting server at port 8080")

	http.ListenAndServe(":8080", r)
}
