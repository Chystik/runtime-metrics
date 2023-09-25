SHELL := /bin/bash
.PHONY: test-one test-all dep test race lint gen cover statictest

dep:
	go mod download
	go mod tidy

test:
	go test ./...

race:
	go test -v -race ./...

lint:
	/home/user/go/bin/golangci-lint run

# use: make test-one iter=9
agent-bin = ./cmd/agent/agent
server-bin = ./cmd/server/server
test-one:
	go build -o $(agent-bin) ./cmd/agent
	go build -o $(server-bin) ./cmd/server
	./metricstest -test.v \
	-test.run=^TestIteration$(iter)[A-B]*$$ \
	-agent-binary-path=$(agent-bin) \
	-binary-path=$(server-bin) \
	-server-port=8080 -source-path=. \
	-file-storage-path=./data.json \
	-database-dsn='postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable'
	rm $(agent-bin) $(server-bin)

# use: make test-all iter=9
ifeq ($(shell test $(iter) > 9; echo $$?),0)
 $(eval t := $$$(iter))
 r := $(subst $(iter),,$(t))
 reg='([1-9]|1[0-$(r)])[A-B]*'
else
 reg=[1-$(iter)][A-B]*
endif
test-all:
	go build -o $(agent-bin) ./cmd/agent
	go build -o $(server-bin) ./cmd/server
	./metricstest -test.v \
	-test.run=^TestIteration$(reg)$$ \
	-agent-binary-path=$(agent-bin) \
	-binary-path=$(server-bin) \
	-server-port=8080 \
	-source-path=. \
	-file-storage-path=./data.json \
	-database-dsn='postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable'
	rm $(agent-bin) $(server-bin)

gen:
	go generate ./...

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o=coverage.html
	rm coverage.out

statictest:
	go vet -vettool=statictest ./...