package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	handlebarsJSURL = "https://cdnjs.cloudflare.com/ajax/libs/handlebars.js/4.7.8/handlebars.min.js"
	outputFile      = "pkg/jsruntime/handlebarsjs.gen.go"
)

func main() {
	resp, err := http.Get(handlebarsJSURL)
	if err != nil {
		log.Fatal("Error fetching HandlebarsJS:", err)
	}
	defer resp.Body.Close()

	jsCode, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	escapedJSCode := strings.ReplaceAll(string(jsCode), "`", "` + \"`\" + `")

	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	content := fmt.Sprintf("package jsruntime\n\n// Generated code, DO NOT EDIT.\n\nconst HandlebarsJS = `%s`\n", escapedJSCode)

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal("Error writing to file:", err)
	}

}
