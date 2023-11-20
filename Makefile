SHELL = /bin/bash

.PHONY: dep
dep:
	go mod download
	go mod tidy

.PHONY: test
test:
	go test ./...

.PHONY: race
race:
	go test -v -race ./...

.PHONY: lint
lint:
	/home/user/go/bin/golangci-lint run

.PHONY: test-one
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
	-key="ssss" \
	-database-dsn='postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable'
	rm $(agent-bin) $(server-bin)

.PHONY: test-all
# use: make test-all iter=9
ifeq ($(shell test $(iter) -gt 9; echo $$?),0)
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
	-key="ssss" \
	-database-dsn='postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable'
	rm $(agent-bin) $(server-bin)

.PHONY: gen
gen:
	go generate ./...

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o=coverage.html
	rm coverage.out

.PHONY: statictest
statictest:
	go vet -vettool=statictest ./...

.PHONY: dev-up
dev-up:
	docker-compose -f=docker-compose.dev.yml --env-file=.env.dev up -d

.PHONY: dev-down
dev-down:
	docker-compose -f=docker-compose.dev.yml --env-file=.env.dev down --rmi local