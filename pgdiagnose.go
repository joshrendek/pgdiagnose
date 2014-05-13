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
	Url string `json:"url" binding:"required"`
}

func create(params JobParams, db *sql.DB) (int, string) {
	var id string
	row := db.QueryRow("INSERT INTO results DEFAULT VALUES returning id")
	err := row.Scan(&id)
	if err != nil {
		log.Print("%v", err)
		return 500, "error"
	}
	fmt.Println("id: %v", id)
	return 200, "foo: " + id
}

func setupDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
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
