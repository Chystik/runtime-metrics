.PHONY: test-local dep test race lint gen cover

dep:
	go mod download
	go mod tidy

test:
	go test ./...

race:
	go test -v -race ./...

lint:
	/home/user/go/bin/golangci-lint run

# use: make test-local iter=2
test-local:
	go build -o ./cmd/agent/agent ./cmd/agent
	go build -o ./cmd/server/server ./cmd/server
	./metricstest -test.v -test.run=^TestIteration$(iter)$$ -agent-binary-path=./cmd/agent/agent -binary-path=./cmd/server/server -server-port=8080 -source-path=./
	rm ./cmd/agent/agent ./cmd/server/server

gen:
	go generate ./...

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o=coverage.html
	rm coverage.out