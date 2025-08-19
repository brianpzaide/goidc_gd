package main

import (
	"flag"
	"goidc_gd/models/store"
	"goidc_gd/models/store/sqlite"
	"html/template"
	"log"
	"os"
)

const (
	MODELS_DSN   = "./goidc_gd.db"
	SESSIONS_DSN = "./goidc_gd_sessions.db"
)

type application struct {
	logger         *log.Logger
	models         store.Models
	sessionManager *sqlite.SessionManagerApp
	clientID       string
	clientSecret   string
}

var (
	dsn          string
	port         int
	homePageTmpl *template.Template
	loginPage    []byte
	clientId     string
	clientSecret string
)

func main() {

	flag.IntVar(&port, "port", 4000, "API server port")
	flag.StringVar(&dsn, "dsn", MODELS_DSN, "SQLITE3 DSN")
	flag.StringVar(&clientId, "cid", os.Getenv("CLIENT_ID"), "client ID")
	flag.StringVar(&clientSecret, "csecret", os.Getenv("CLIENT_SECRET"), "Client Secret")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	if clientId == "" || clientSecret == "" {
		logger.Fatal("clientid or client secret must not be empty")
	}

	var (
		m   store.Models
		sm  *sqlite.SessionManagerApp
		err error
	)

	m, err = store.NewModel(dsn)
	if err != nil {
		logger.Fatal(err)
	}

	defer m.Close()

	sm, err = store.NewSessionManager(SESSIONS_DSN)
	if err != nil {
		logger.Fatal(err)
	}
	defer sm.Close()

	homePageTmpl, err = template.ParseFiles("./ui/html/index.html")
	if err != nil {
		logger.Fatal(err)
	}

	loginPage, err = os.ReadFile("./ui/html/login.html")
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		logger:         logger,
		models:         m,
		sessionManager: sm,
	}

	err = app.serve()
	if err != nil {
		logger.Fatal(err)
	}
}
