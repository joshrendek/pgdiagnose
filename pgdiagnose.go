package main

import (
	_ "fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

type JobParams struct {
	Url string `json:"url" binding:"required"`
}

func create(params JobParams) (int, string) {
	return 200, params.Url
}

func main() {
	m := martini.Classic()
	m.Post("/create", binding.Json(JobParams{}), create)
	m.Run()
}
