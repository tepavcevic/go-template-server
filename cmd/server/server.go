package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/tepavcevic/go-template-server/controllers"
	"github.com/tepavcevic/go-template-server/migrations"
	"github.com/tepavcevic/go-template-server/models"
	"github.com/tepavcevic/go-template-server/templates"
	"github.com/tepavcevic/go-template-server/views"
)

type config struct {
	PSQL   models.PostgresConfig
	SMTP   models.SMTPConfig
	Server struct {
		Address string
	}
	CSRF struct {
		Key    string
		Secure bool
	}
}

func loadEnvConfig() (config, error) {
	var cfg config

	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DBName:   os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	if cfg.PSQL.Host == "" && cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no PSQL config provided")
	}

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	csrfKeyStr := os.Getenv("CSRF_SECURE")
	cfg.CSRF.Secure = csrfKeyStr == "true"

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	err = run(cfg)
	if err != nil {
		panic(err)
	}
}

func run(cfg config) error {
	// Database setup
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		return err
	}

	// Services setup
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	passwordResetService := &models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)
	galleryService := models.GalleryService{
		DB: db,
	}

	// Middleware setup
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}
	csrfMiddleware := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	// Controllers setup
	userC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: passwordResetService,
		EmailService:         emailService,
	}
	galleryC := controllers.Galleries{
		GalleryService: &galleryService,
	}
	userC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "tailwind.gohtml",
	))
	userC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "tailwind.gohtml",
	))
	userC.Templates.ForgotPassword = views.Must(views.ParseFS(
		templates.FS,
		"forgot-password.gohtml", "tailwind.gohtml",
	))
	userC.Templates.CheckYourEmail = views.Must(views.ParseFS(
		templates.FS,
		"check-your-email.gohtml", "tailwind.gohtml",
	))
	userC.Templates.ResetPassword = views.Must(views.ParseFS(
		templates.FS,
		"reset-password.gohtml", "tailwind.gohtml",
	))
	galleryC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"galleries/new.gohtml", "tailwind.gohtml",
	))
	galleryC.Templates.Edit = views.Must(views.ParseFS(
		templates.FS,
		"galleries/edit.gohtml", "tailwind.gohtml",
	))
	galleryC.Templates.Index = views.Must(views.ParseFS(
		templates.FS,
		"galleries/index.gohtml", "tailwind.gohtml",
	))
	galleryC.Templates.Show = views.Must(views.ParseFS(
		templates.FS,
		"galleries/show.gohtml", "tailwind.gohtml",
	))

	// Router and routes setup
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(csrfMiddleware)
	r.Use(umw.SetUser)

	r.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(
		templates.FS,
		"home.gohtml", "tailwind.gohtml",
	))))
	r.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(
		templates.FS,
		"contact.gohtml", "tailwind.gohtml",
	))))
	r.Get("/faq", controllers.FAQ(views.Must(views.ParseFS(
		templates.FS,
		"faq.gohtml", "tailwind.gohtml",
	))))
	r.Get("/signup", userC.New)
	r.Post("/users", userC.Create)
	r.Get("/signin", userC.SignIn)
	r.Post("/signin", userC.ProcessSignIn)
	r.Post("/signout", userC.ProcessSignOut)
	r.Get("/forgot-pw", userC.ForgotPassword)
	r.Post("/forgot-pw", userC.ProcessForgotPassword)
	r.Get("/reset-pw", userC.ResetPassword)
	r.Post("/reset-pw", userC.ProcessResetPassword)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", userC.CurrentUser)
	})
	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleryC.Show)
		r.Get("/{id}/images/{filename}", galleryC.Image)

		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/", galleryC.Index)
			r.Get("/new", galleryC.New)
			r.Post("/", galleryC.Create)
			r.Get("/{id}/edit", galleryC.Edit)
			r.Post("/{id}", galleryC.Update)
			r.Post("/{id}/delete", galleryC.Delete)
			r.Post("/{id}/images/{filename}/delete", galleryC.DeleteImage)
			r.Post("/{id}/images", galleryC.UploadImage)
		})
	})
	assetsHandler := http.FileServer(http.Dir("assets"))
	r.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Println("Starting server at port", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, r)
}
