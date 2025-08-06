package main

import (
	"log"

	"github.com/DanielLiu1123/gencoder/pkg/cmd"
	"github.com/DanielLiu1123/gencoder/pkg/model"

	// drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
)

var version = "0.1.1"

func main() {

	log.SetFlags(0)

	buildInfo := &model.BuildInfo{
		Version: version,
	}

	err := cmd.NewCmdRoot(buildInfo).Execute()
	if err != nil {
		log.Fatal(err)
	}
}
