package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"

	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type JobParams struct {
	URL string `json:"url" binding:"required"`
}

func createJob(db *sql.DB) (id string, err error) {
	row := db.QueryRow("INSERT INTO results DEFAULT VALUES returning id")
	err = row.Scan(&id)
	if err != nil {
		log.Print("%v", err)
		return "", err
	}
	fmt.Println("new job id: %v", id)
	return id, nil
}

func create(params JobParams, db *sql.DB) (int, string) {
	id, err := createJob(db)
	if err != nil {
		log.Print("%v", err)
		return 500, "error"
	}
	return 200, "foo: " + id
}

func setupDB() *sql.DB {
	connstring := os.Getenv("DATABASE_URL")
	if connstring == "" {
		connstring = "dbname=pgdiagnose sslmode=disable"
	}

	db, err := sql.Open("postgres", connstring)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("select 1")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	m := martini.Classic()

	m.Map(setupDB())
	m.Post("/create", binding.Json(JobParams{}), create)
	m.Run()
}
