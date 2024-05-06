package main

import (
	"fmt"

	"github.com/nagy-gergely/api-test/cmd/api"
	"github.com/nagy-gergely/api-test/db"
)

func main() {
	db := db.NewSqlite()

	server := api.NewAPIServer(":8080", db)

	if err := server.Start(); err != nil {
		fmt.Println(err)
	}
}
