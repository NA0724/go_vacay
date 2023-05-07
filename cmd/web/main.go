package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"Go_Vacay/internal/config"
	"Go_Vacay/internal/handlers"
	"Go_Vacay/internal/models"
	"Go_Vacay/internal/renderers"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var session *scs.SessionManager
var app config.AppConfig

// main is the main function
func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Staring application on port %s \n", portNumber)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: Routes(&app),
	}
	err = serve.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	gob.Register(models.Reservation{})
	gob.Register(models.Registration{})
	gob.Register(models.Login{})
	//set to true if production environment
	app.InProd = false

	//initialise session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProd

	app.Session = session

	tempCache, err := renderers.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err)
		return err
	}

	app.TemplateCache = tempCache
	app.UseCache = false // false for development mode, true for prod or qa mode

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	renderers.NewTemplates(&app)
	return nil
}
