package main

import (
	"fmt"
	"github.com/will/pgdiagnose"
	"os"
)

func main() {
	connstring := "dbname=will sslmode=disable"
	if len(os.Args) > 1 {
		connstring = os.Args[1]
	}

	checks, _ := pgdiagnose.CheckAll(connstring)
	json, _ := pgdiagnose.PrettyJSON(checks)
	fmt.Println(json)
}
