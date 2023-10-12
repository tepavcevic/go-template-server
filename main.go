package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/tepavcevic/go-template-server/controllers"
	"github.com/tepavcevic/go-template-server/migrations"
	"github.com/tepavcevic/go-template-server/models"
	"github.com/tepavcevic/go-template-server/templates"
	"github.com/tepavcevic/go-template-server/views"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	tpl := views.Must(views.ParseFS(
		templates.FS,
		"home.gohtml", "tailwind.gohtml",
	))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(
		templates.FS,
		"contact.gohtml", "tailwind.gohtml",
	))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(
		templates.FS,
		"faq.gohtml", "tailwind.gohtml",
	))
	r.Get("/faq", controllers.FAQ(tpl))

	cfg := models.DefaultPostgresConfig()
	fmt.Println(cfg)

	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	userC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	userC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	))
	userC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "tailwind.gohtml",
	))
	r.Get("/signup", userC.New)
	r.Post("/users", userC.Create)
	r.Get("/signin", userC.SignIn)
	r.Post("/signin", userC.ProcessSignIn)
	r.Post("/signout", userC.ProcessSignOut)
	r.Get("/users/me", userC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfMiddleware := csrf.Protect([]byte("qwertyuiop34asdfghjkle6754321azxcvbn"), csrf.Secure(false))

	fmt.Println("Starting server at port 8080")
	http.ListenAndServe(":8080", csrfMiddleware(umw.SetUser(r)))
}
