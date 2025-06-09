.PHONY: api-hotreload lint

# Run the backend API with hot-reload, make sure you installed `air`
api-hotreload:
	air --build.cmd "go build -o bin/api cmd/api/main.go" --build.bin "./bin/api"

lint:
	golangci-lint run ./...
