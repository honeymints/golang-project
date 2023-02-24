package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"todolist.net/internal/data"
	"todolist.net/internal/mailer"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

// Define an application struct to hold the dependencies for our HTTP handlers, helpers,
// and middleware. At the moment this only contains a copy of the config struct and a
// logger, but it will grow to include a lot more as our build progresses.
type application struct {
	config config
	logger *log.Logger
	models data.Models
	mailer mailer.Mailer
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 3000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp-mail.outlook.com", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "211307@astanait.edu.kz", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "gudron5027962003", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "211307@astanait.edu.kz", "SMTP sender")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db, err := sql.Open("postgres", "postgres://admin:gudron@hostname:5432/todolistsslmode=require")

	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// Also log a message to say that the connection pool has been successfully
	// established.
	logger.Print("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}
	// Use the httprouter instance returned by app.routes() as the server handler.
	srv := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err = srv.ListenAndServe()

	logger.Fatal(err)

}
