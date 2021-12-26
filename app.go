package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed templates
var templatesFS embed.FS

type Application struct {
	templates fs.FS
	mailer    Mailer
	logger    *log.Logger
}

func NewApplication(mailer Mailer, logger* log.Logger) *Application {
	var templates fs.FS
	if strings.HasPrefix(os.Getenv("ENV"), "dev") {
		// reload templates from filesystem if var ENV starts with "dev"
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		templates = os.DirFS("./templates/")
	} else {
		// else use the embedded templates dir
		templates, _ = fs.Sub(templatesFS, "templates")
	}

	app := Application{
		templates: templates,
		mailer:    mailer,
		logger:    logger,
	}

	return &app
}

func (app *Application) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.NotFound(app.notFoundResponse)
	r.MethodNotAllowed(app.methodNotAllowedResponse)

	r.Get("/", app.HandleIndex)
	r.Post("/contact", app.HandleContact)

	return r
}

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFS(app.templates, "index.html.tmpl")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) HandleContact(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	email := r.PostFormValue("email")
	if email == "" {
		http.Redirect(w, r, "/", 303)
		return
	}

	from := "info@shallowbrooksoftware.com"
	to := "info@shallowbrooksoftware.com"
	subject := "Business Inquiry from Website!"
	body := fmt.Sprintf("Someone wants to get in touch:\n%s", email)

	err = app.mailer.SendMail(from, from, to, to, subject, body)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	ts, err := template.ParseFS(app.templates, "thanks.partial.tmpl")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, tmpl string) {
	// attempt to parse error template
	ts, err := template.ParseFS(app.templates, tmpl)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}

	// render template to a temp buffer
	var buf bytes.Buffer
	err = ts.Execute(&buf, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}

	// write the status and error page
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, 404, "404.html.tmpl")
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, 405, "405.html.tmpl")
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// skip 2 frames to identify original caller
	app.logger.Output(2, err.Error())
	app.errorResponse(w, r, 500, "500.html.tmpl")
}
