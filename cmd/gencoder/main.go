package main

import (
	"github.com/DanielLiu1123/gencoder/pkg/cmd"
	"log"

	// drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/sijms/go-ora/v2"
)

func main() {

	log.SetFlags(0)

	err := cmd.NewCmdRoot().Execute()
	if err != nil {
		log.Fatal(err)
	}
}