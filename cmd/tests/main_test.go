package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"todolist.net/internal/data"
)

var user *data.UserModel

const (
	dbDriver = "postgres"
	dbSource = "postgres://gatgwnrc:4ATZR9hsHW4jnQRKgMvhM9310IoCvZDv@trumpet.db.elephantsql.com/gatgwnrc"
)

func TestMain(m *testing.M) {
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("couldn't connect to database")
	}
	log.Println("TestMain is executing")
	code := m.Run()

	db.Close()
	os.Exit(code)
}

/*
	 func TestMain(m *testing.M) {
		db, err := sql.Open(dbDriver, dbSource)
		if err != nil {
			log.Fatal("couldn't connect to database")
		}
		log.Println("TestMain is executing")
		code := m.Run()

		db.Close()
		os.Exit(code)
	}
*/
func TestSomething(t *testing.T) {
	t.Log("works fine")
}
