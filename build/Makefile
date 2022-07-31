.PHONY: test
export GO111MODULE ?= on
TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
GIT_VERSION = $(shell git describe --always --abbrev=7 --dirty)
SHELL = /bin/sh

export COV_FOLDER?=coverage
export APP_NAME?=warehouse
export COVERAGE_PROFILE?=$(COV_FOLDER)/cover.out
export COVERAGE_HTML?=$(COV_FOLDER)/index.html
export PROJECT_DIR ?= $(PWD)
GORACE="log_path=./report_race"
LINT_TOOL ?= ./bin/golangci-lint

#####################################
# Help
#####################################
help: ## Display help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


#####################################
# Package dependencies
#####################################
vendor: dep ## Run vendor tidy
	go mod vendor

dep: ## Fetch dependencies
	go install .
	go mod tidy
	go mod verify


#####################################
# Tests (coverage benchmarcking)
#####################################

test: dep vendor ## Run tests
	go test ./...

# Generate test coverage
test-cov:     ## Run test coverage and generate html report
	rm -fr $(COV_FOLDER)
	mkdir $(COV_FOLDER)
	go list -f '{{if gt (len .TestGoFiles) 0}}"go test -covermode count -coverprofile {{.Name}}.coverprofile -coverpkg ./... {{.ImportPath}}"{{end}}' ./... | xargs -I {} bash -c {}
	echo "mode: count" > $(COV_FOLDER)/cover.out
	grep -h -v "^mode:" *.coverprofile >> $(COV_FOLDER)"/cover.out"
	rm *.coverprofile
	go tool cover -html=coverage/cover.out -o=$(COV_FOLDER)/cover.html

race: dep ## Run race detector
	@GORACE="log_path=./report_race" go test -race -short ./...


#####################################
# Lint&format
#####################################

fmt:  ## Execute go fmt
	go fmt ./...

install-lint: ## Install the dependencies to run the linter
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.47.2

lint:
	@echo "Running golangci-lint"
	$(LINT_TOOL) run -c build/.golangci.yml


#####################################
# pack & run-server & stop-server
#####################################
pack:
	docker build . -f ./build/Dockerfile -t  $(APP_NAME)

run-server: pack
	cd build && docker-compose up -d warehouse

stop-server:
	cd build && docker-compose stop warehouse

#####################################
# database migration
#####################################
migrate:
	liquibase --username=postgres --password=password --url jdbc:postgresql://localhost:5432/warehouse --changeLogFile=liquibase/changelog/master.yml --logLevel=info update


#####################################
# Mocks
#####################################

mocks-generate: ## update mocks
	mockery --dir ./app --all --disable-version-string


mocks-install: ## install mockery (once)
	go install github.com/vektra/mockery/v2@v2.10.0


#####################################
# Build&run&clean
#####################################
build: dep ## Build executales
	go build -o $(APP_NAME) .



run: build  ## Run executable
	./$(APP_NAME)

clean:  ## Clean
	rm $(APP_NAME)
	rm -rf  $(COV_FOLDER)
