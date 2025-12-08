#!/usr/bin/make -f


format-tools:
	go install mvdan.cc/gofumpt@v0.3.1
	go install github.com/client9/misspell/cmd/misspell@v0.3.4
	go install golang.org/x/tools/cmd/goimports@latest

lint: format-tools
	golangci-lint run --tests=false
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*_test.go" | xargs gofumpt -d -s


format: format-tools
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofumpt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs goimports -w -local github.com/thesixnetwork/lbb-sdk-go