package main

import (
	"encoding/gob"
	"fmt"
	"github.com/easymomo/go-bookings/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/easymomo/go-bookings/internal/config"
	"github.com/easymomo/go-bookings/internal/handlers"
	"github.com/easymomo/go-bookings/internal/render"

	"github.com/alexedwards/scs/v2"
)

// Variables declared here (before the main function) are available to every function in the main package
const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	// Things we can store in the session, primitives and string can be added by default but we need
	// to tell the application about structs we created ourselves
	gob.Register(models.Reservation{})

	// Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	render.NewTemplates(&app)

	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	fmt.Println(fmt.Printf("Starting the server on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
