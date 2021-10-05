package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/sbs-website/mail"
)

type Application struct {
	mailer mail.Mailer
}

func (app *Application) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("templates/index.html.tmpl")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err)
		return
	}
}

func (app *Application) HandleContact(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	email := r.PostFormValue("email")
	if email == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	from := "info@shallowbrooksoftware.com"
	to := "info@shallowbrooksoftware.com"
	subject := "Business Inquiry from Website!"
	body := fmt.Sprintf("Someone wants to get in touch:\n%s", email)
	err = app.mailer.SendMail(from, from, to, to, subject, body)
	if err != nil {
		log.Println(err)
	}

	// TODO: flash some sort of "we'll be in touch" message
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	// use SendGridMailer in prod else default to LogMailer
	var mailer mail.Mailer
	if os.Getenv("ENV") == "production" {
		sendGridAPIKey := os.Getenv("SENDGRID_API_KEY")
		if sendGridAPIKey == "" {
			log.Fatal("Missing required env var: SENDGRID_API_KEY")
		}
		mailer = mail.NewSendGridMailer(sendGridAPIKey)
	} else {
		mailer = mail.NewLogMailer()
	}

	app := &Application{
		mailer: mailer,
	}

	router := httprouter.New()
	router.HandlerFunc("GET", "/", app.HandleIndex)
	router.HandlerFunc("POST", "/contact", app.HandleContact)
	router.Handler("GET", "/metrics", promhttp.Handler())
	router.ServeFiles("/posts/*filepath", http.Dir("./posts"))
	router.ServeFiles("/static/*filepath", http.Dir("./static"))

	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
