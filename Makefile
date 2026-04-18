.PHONY: all
 all: help

GIT_TAG := $(shell git rev-parse --short HEAD)
SERVICE_IMAGE := docker.io/nghiango1/interingo-service
NVIM_IMAGE := docker.io/nghiango1/interingo

### Build command

.PHONY: build
build: clean website/dist interingo-build interingo-service-build lsp-build # Build all the code: Website front-end, REPL, backend-service, lsp

.PHONY: interingo-build
interingo-build: embed-content-clean # Build Interingo REPL and API only, doesn't contain front-end
	mkdir -p dist
	go build -o dist/interingo cmd/interingo/main.go

.PHONY: interingo-service-build
interingo-service-build: embed-content # Build Interingo Backend services with front-end packed
	go build -o dist/interingo cmd/interingo/main.go

.PHONY: lsp-build
lsp-build: # Build go binary file
	mkdir -p dist
	go build -o dist/interingo-lsp cmd/interingo-lsp/main.go

### Front-end

website/dist: # Build website static file, output into website/dist 
	cd website/ && npm install && npm run build

.PHONY: website/dist-force
website/dist-force: # Force rebuild website static file, output into website/dist 
	cd website/ && npm install && npm run build

.PHONY: embed-content
embed-content: website/dist embed-content-clean # Put built website into embed directory for go compile
	cp -r website/dist/ pkg/server/content

### Container deploy helper
.PHONY: docker-build
docker-build: # Build the container image for services hosting
	docker build -f docker/service.Dockerfile . -t $(SERVICE_IMAGE):latest
	docker tag $(NVIM_IMAGE):latest $(NVIM_IMAGE):$(GIT_TAG)

.PHONY: docker-push
docker-push: # Push the services hosting image into docker.io
	docker push $(NVIM_IMAGE):latest
	docker push $(NVIM_IMAGE):$(GIT_TAG)

.PHONY: docker-nvim-build
docker-nvim-build: # Build the image for nvim showcase
	docker build -f docker/nvim.Dockerfile . -t $(SERVICE_IMAGE):latest
	docker tag $(SERVICE_IMAGE):latest $(SERVICE_IMAGE):$(GIT_TAG)

.PHONY: docker-nvim-push
docker-nvim-push: # Push the image for nvim showcase into docker.io
	docker push $(SERVICE_IMAGE):latest
	docker push $(SERVICE_IMAGE):$(GIT_TAG)


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

### Helper / Clean-up

.PHONY: website-clean
website-clean: # Clean up website built
	rm -rf website/dist/

.PHONY: embed-content-clean
embed-content-clean: # Clean up old webpage content, 
	rm -rf pkg/server/content/**

.PHONY: go-clean
go-clean: # Clean up go built files 
	rm -rf dist/

.PHONY: clean
clean: website-clean embed-content-clean go-clean # Clean up all built file

.PHONY: help
help: # Show this help
	@cat Makefile | \
		grep -E '^[^.:[:space:]]+:|[#]##' | \
		sed -E 's/:[^#]*#([^:]+)$$/: #:\1/' | \
		sed -E 's/([^.:[:space:]]+):([^#]*#(.+))?.*/    make \1\3/' | \
		sed -E 's/[#][#]# *(.+)/# \1/' | \
		column -ts: -L -W2
