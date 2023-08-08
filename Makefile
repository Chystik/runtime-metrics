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
agent-bin = ./cmd/agent/agent
server-bin = ./cmd/server/server
test-local:
	go build -o $(agent-bin) ./cmd/agent
	go build -o $(server-bin) ./cmd/server
	./metricstest -test.v -test.run=^TestIteration$(iter)$$ -agent-binary-path=$(agent-bin) -binary-path=$(server-bin) -server-port=8080 -source-path=.
	rm $(agent-bin) $(server-bin)

gen:
	go generate ./...

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o=coverage.html
	rm coverage.out