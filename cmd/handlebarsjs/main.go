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
	handlebarsJSURL    = "https://cdnjs.cloudflare.com/ajax/libs/handlebars.js/4.7.8/handlebars.min.js"
	handlebarsJSOutput = "pkg/jsruntime/handlebarsjs.gen.go"

	helperJSFile = "pkg/jsruntime/helper.js"
	helperOutput = "pkg/jsruntime/helper.gen.go"
)

func main() {

	// Generate handlebarsjs.gen.go
	genHandlebarJS()

	// Generate helper.gen.go
	genHelper()

}

func genHelper() {
	jsCode, err := os.ReadFile(helperJSFile)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	escapedJSCode := strings.ReplaceAll(string(jsCode), "`", "` + \"`\" + `")

	file, err := os.Create(helperOutput)
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	content := fmt.Sprintf("package jsruntime\n\n// Generated file, DO NOT EDIT.\n\nconst HelperJS = `%s`\n", escapedJSCode)

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal("Error writing to file:", err)
	}

}

func genHandlebarJS() {
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

	file, err := os.Create(handlebarsJSOutput)
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	content := fmt.Sprintf("package jsruntime\n\n// Generated file, DO NOT EDIT.\n\nconst HandlebarsJS = `%s`\n", escapedJSCode)

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatal("Error writing to file:", err)
	}
}
