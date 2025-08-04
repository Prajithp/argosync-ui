.PHONY: all build clean mod generate help

# custom define
PROJECT := argosync
MAINFILE := main.go

all: build

mod: ## Get the dependencies
	go mod download

generate: ## generate the static assets
	go generate ./...

build: mod ## Build the binary file
	go build -v -o build/bin/$(PROJECT) $(MAINFILE)


clean: ## Remove previous build
	@rm -rf ./build

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
