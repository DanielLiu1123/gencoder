
init :
	git config core.hooksPath .githooks
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.3.1

.PHONY : init
