package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/sangketkit01/bookings/internal/config"
	"github.com/sangketkit01/bookings/internal/driver"
	"github.com/sangketkit01/bookings/internal/handlers"
	"github.com/sangketkit01/bookings/internal/helpers"
	"github.com/sangketkit01/bookings/internal/models"
	"github.com/sangketkit01/bookings/internal/render"
)

const portNumber = ":8216"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil{
		log.Fatalln(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)
	listenForMail()

	server := &http.Server{
		Addr:    portNumber,
		Handler: route(&app),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() (*driver.DB,error){
	// What am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})
	gob.Register(map[string]int{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	//Change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t",log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout,"ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=0627457454New")
	if err != nil{
		log.Fatal("Cannot connect to the database! Dying...")
	}

	log.Println("Connected to the database!")


	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalln("cannot create template cache", err)
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepository(&app,db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	fmt.Println(fmt.Sprintf("Listening on port %s", portNumber))

	return db,nil
}


