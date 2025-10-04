init :
	git config core.hooksPath .githooks
	# https://golangci-lint.run/docs/welcome/install/#binaries
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.5.0
