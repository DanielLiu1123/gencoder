
init :
	git config core.hooksPath .githooks
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

# update handlebars
hb :
	go run cmd/handlebarsjs/main.go

.PHONY : init
