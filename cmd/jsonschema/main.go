package main

import (
	"github.com/DanielLiu1123/gencoder/pkg/model"
	"github.com/DanielLiu1123/gencoder/pkg/util"
	"github.com/invopop/jsonschema"
	"log"
)

const jsonschemaFile = "schema.json"

//go:generate sh -c "cd ../../ && go run cmd/jsonschema/main.go"
func main() {
	schema := jsonschema.Reflect(&model.Config{})
	schema.ID = "https://github.com/DanielLiu1123/gencoder/tree/main/pkg/model/table.go"

	err := util.WriteFile(jsonschemaFile, []byte(util.ToJson(schema)))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("JSON schema generated at", jsonschemaFile)
}
