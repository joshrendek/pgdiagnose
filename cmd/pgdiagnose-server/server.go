package main

import (
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/binding"
	"github.com/will/pgdiagnose"
	"log"
	"net/http"
	"os"
)

type JobParams struct {
	URL string `json:"url" binding:"required"`
}

func getResultJSON(id string, db *sql.DB) (json string, err error) {
	row := db.QueryRow("SELECT row_to_json(results, true) FROM results WHERE id = $1", id)
	err = row.Scan(&json)
	if err != nil {
		log.Print("%v", err)
		return "", err
	}
	return json, nil
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

	return id, nil
}

func create(params JobParams, db *sql.DB) (int, string) {
	id, err := createJob(db, params)
	if err != nil {
		log.Print("%v", err)
		return 500, "error"
	}

	json, err2 := getResultJSON(id, db)
	if err2 != nil {
		log.Print("%v", err2)
		return 500, "error"
	}

	return 201, json
}

func getReport(params martini.Params, db *sql.DB) (int, string) {
	json, err := getResultJSON(params["id"], db)
	if err != nil {
		return 404, ""
	}

	return 200, json
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

	if martini.Env == "production" {
		m.Use(func(res http.ResponseWriter, req *http.Request) {
			if req.Header.Get("X-FORWARDED-PROTO") != "https" {
				fmt.Println("not https: ", req.Header.Get("X-FORWARDED-PROTO"))
				res.WriteHeader(http.StatusUnauthorized)
			}
		})
	}
	m.Map(setupDB())
	m.Post("/reports", binding.Json(JobParams{}), create)
	m.Get("/reports/:id", getReport)
	m.Run()
}
