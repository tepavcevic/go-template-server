package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func handlerHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func handlerContact(w http.ResponseWriter, r *http.Request) {
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
	r.Get("/", handlerHome)
	r.Get("/contact", handlerContact)
	r.Get("/faq", handlerFAQ)

	fmt.Println("Starting server at port 8080")

	http.ListenAndServe(":8080", r)
}
