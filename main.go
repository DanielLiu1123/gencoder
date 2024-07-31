package main

import (
	"fmt"
	"github.com/DanielLiu1123/gencoder/cmd"
	"github.com/aymerick/raymond"
	_ "github.com/aymerick/raymond"

	// drivers
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/sijms/go-ora/v2"
)

func main() {
	cmd.Execute()

	source := `<div class="entry">
	 <h1>{{title}}</h1>
	 <div class="body">
	   {{body}}
	 </div>
	</div>
	`
	ctxList := []map[string]string{
		{
			"title": "My New Post",
			"body":  "This is my first post!",
		},
		{
			"title": "Here is another post",
			"body":  "This is my second post!",
		},
	}

	// parse template
	tpl, err := raymond.Parse(source)
	if err != nil {
		panic(err)
	}

	for _, ctx := range ctxList {
		// render template
		result, err := tpl.Exec(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Print(result)
	}

}
