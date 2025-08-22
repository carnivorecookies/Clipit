all: clipit

clipit: check
	go build ./cmd/clipit

check: fix
	-golangci-lint-v2 run

fix:
	go mod tidy
	gofmt -s -w .
