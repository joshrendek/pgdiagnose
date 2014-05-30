package main

import (
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"github.com/martini-contrib/binding"
	"log"
	"net/http"
	"os"
	"regexp"
)

type JobParams struct {
	URL     string `json:"url" binding:"required"`
	Metrics []struct {
		LoadAvg1m *float64 `json:"load_avg_1m"`
	}
	Plan     string
	App      string
	Database string
}

var validParams = regexp.MustCompile(`\A[a-zA-Z0-9\-_]+\z`)

func (params *JobParams) sanitize() {
	if !validParams.MatchString(params.Plan) {
		params.Plan = ""
	}
	if !validParams.MatchString(params.App) {
		params.App = ""
	}
	if !validParams.MatchString(params.Database) {
		params.Database = ""
	}
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
	params.sanitize()

	plan := GetPlan(params.Plan)

	checks, err := CheckSql(params.URL, plan)
	if err != nil {
		return "", err
	}

	fmt.Println(params.Metrics)

	loadChecks := func() []Check {
		if len(params.Metrics) > 0 {
			return CheckLoad(params.Metrics[0].LoadAvg1m)
		} else {
			return CheckLoad(nil)
		}
	}()
	checks = append(checks, loadChecks...)

	checksJSON, _ := PrettyJSON(checks)

	row := db.QueryRow("INSERT INTO results (app,database,checks) values ($1,$2,$3) returning id", params.App, params.Database, checksJSON)
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

func health(db *sql.DB) (int, string) {
	_, err := db.Exec("select 1")
	if err != nil {
		log.Println(err)
		return 500, "database error"
	}
	return 200, "ok"
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
	m.Get("/health", health)
	m.Run()
}
