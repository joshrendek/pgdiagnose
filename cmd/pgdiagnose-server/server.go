package main

import (
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/binding"
	"github.com/will/pgdiagnose"
	"log"
	"os"
)

type JobParams struct {
	URL string `json:"url" binding:"required"`
}

func createJob(db *sql.DB, params JobParams) (id string, err error) {
	checks := pgdiagnose.CheckAll(params.URL)
	row := db.QueryRow("INSERT INTO results (data) values ($1) returning id", checks)
	err = row.Scan(&id)
	if err != nil {
		log.Print("%v", err)
		return "", err
	}
	fmt.Println("new job id: %v", id)
	return id, nil
}

func create(params JobParams, db *sql.DB) (int, string) {
	id, err := createJob(db, params)
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
