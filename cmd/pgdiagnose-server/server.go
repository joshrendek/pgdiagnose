package main

import (
	"database/sql"
	"encoding/json"
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

type responseWithId struct {
	Id     string
	Checks []pgdiagnose.Check
}

func createJob(db *sql.DB, params JobParams) (id string, err error) {
	checks := pgdiagnose.CheckAll(params.URL)

	checksJSON, _ := pgdiagnose.PrettyJSON(checks)

	row := db.QueryRow("INSERT INTO results (checks) values ($1) returning id", checksJSON)
	err = row.Scan(&id)
	if err != nil {
		log.Print("%v", err)
		return "", err
	}

	fmt.Println("new job id: ", id)

	response := responseWithId{id, checks}
	responseJSON, _ := json.MarshalIndent(response, "", "  ")

	return string(responseJSON), nil
}

func create(params JobParams, db *sql.DB) (int, string) {
	id, err := createJob(db, params)
	if err != nil {
		log.Print("%v", err)
		return 500, "error"

	}
	return 200, id
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
