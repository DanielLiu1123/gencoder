
init :
	git config core.hooksPath .githooks
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.7

.PHONY : init
