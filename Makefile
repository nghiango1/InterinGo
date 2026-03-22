.PHONY: all
 all: help

### Build command

.PHONY: build
build: embed-dist embed-content go-build # Build all the code

.PHONY: build-run
build-run: build run # Build and run the code

.PHONY: embed-dist
embed-dist: # Build website static file, output into website/dist 
	cd website/ && npm install && npm run build

.PHONY: embed-content
embed-content: # Build webpage then output it into embed directory for go compile
	rm -rf pkg/server/content/**
	@if [[ -e website/dist/ ]]; then \
	  cp -r website/dist/ pkg/server/content/; \
	else \
	  touch pkg/server/content/.not_support; \
	fi

.PHONY: go-build
go-build: # Build go binary file
	mkdir -p dist
	go build -o dist/interingo cmd/interingo/main.go

.PHONY: run clean server-clean repl-clean
run: # Run the build file in server mode
	./dist/interingo -s

### Container deploy helper
.PHONY: docker-build
docker-build: # Build the container image
	docker build -f docker/service.Dockerfile . -t docker.io/nghiango1/interingo-service:latest

.PHONY: docker-push
docker-push: # Push the image into docker.io
	docker push docker.io/nghiango1/interingo-service:latest

### Development helper

.PHONY: go-run
go-run: embed-content # Run the code without build step in server mode
	go run ./cmd/interingo/ -s

.PHONY: regression-test
regression-test: # Run the code without build step in server mode
	python test/regressionTesting.py

.PHONY: go-test
go-test: # Go lang Unit test 
	mkdir -p ./dist
	go test -cover -coverprofile=./dist/coverage.out -coverpkg=./cmd/...,./pkg/...  ./cmd/... ./pkg/...
	go tool cover -func=./dist/coverage.out
	go tool cover -html=./dist/coverage.out -o ./dist/coverage.html
	xdg-open ./dist/coverage.html

### Helper

.PHONY: help
help: # Show this help
	@cat Makefile | \
		grep -E '^[^.:[:space:]]+:|[#]##' | \
		sed -E 's/:[^#]*#([^:]+)$$/: #:\1/' | \
		sed -E 's/([^.:[:space:]]+):([^#]*#(.+))?.*/    make \1\3/' | \
		sed -E 's/[#][#]# *(.+)/# \1/' | \
		column -ts: -L -W2
