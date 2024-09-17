package main

import (
	"github.com/DanielLiu1123/gencoder/pkg/cmd"
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"log"
	"time"

	// drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/sijms/go-ora/v2"
)

var version = "0.0.1"

func main() {

	log.SetFlags(0)

	buildInfo := &model.BuildInfo{
		Version:   version,
		BuildTime: time.Now(),
	}

	err := cmd.NewCmdRoot(buildInfo).Execute()
	if err != nil {
		log.Fatal(err)
	}
}
