package controllers

import (
	"html/template"
	"net/http"
)

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is there a free version",
			Answer:   "Yes! We offer a free trial for 30 days on any period plans.",
		}, {
			Question: "What are Your support hours?",
			Answer:   "We have support staff answering emails 24/7, though response time may be slower on weekends.",
		},
		{
			Question: "How do I contact support?",
			Answer:   `Email us - <a href="mailto:support@example.com">support@example.com</a>.`,
		},
		{
			Question: "Where is your office?",
			Answer:   "Our support team is based in Sekovici, BiH.",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, questions)
	}
}
