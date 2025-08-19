package main

import (
	"flag"
	"goidc_gd/models/store"
	"html/template"
	"log"
	"os"
)

const SQLITE_DSN = "./goidc_gd.db"

type application struct {
	logger *log.Logger
	models store.Models
}

var (
	dsn          string
	port         int
	homePageTmpl *template.Template
	loginPage    []byte
)

func main() {

	flag.IntVar(&port, "port", 4000, "API server port")
	flag.StringVar(&dsn, "db-dsn", SQLITE_DSN, "SQLITE3 DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	var (
		m   store.Models
		err error
	)

	m, err = store.New(dsn)
	if err != nil {
		logger.Fatal(err)
	}

	defer m.Close()

	homePageTmpl, err = template.ParseFiles("./ui/html/index.html")
	if err != nil {
		logger.Fatal(err)
	}

	loginPage, err = os.ReadFile("./ui/html/login.html")
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		logger: logger,
		models: m,
	}

	err = app.serve()
	if err != nil {
		logger.Fatal(err)
	}
}
