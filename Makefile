TAILWIND_CLI := tailwindcss-linux-x64

.PHONY: all build build-run tailwind-build tailwind-watch templ-build templ-watch help go-build run clean server-clean repl-clean

 all: help

### Build command

build: tailwind-build templ-build go-build # Build all the code

build-run: build run # Build and run the code

tailwind-build: # Build tailwind css file, output file server/assets/stylesheet.css
	$(TAILWIND_CLI) -c tailwind.config.js -o server/assets/stylesheet.css

templ-build: # Build/rebuild all `templ` templates files
	templ generate

go-build: # Build go binary file
	go build .

run: # Run the build file in server mode
	./interingo -s

### Development helper

tailwind-watch: # Start tailwind in watch mode - looking for change and rebuild, output file server/assets/stylesheet.css
	$(TAILWIND_CLI) -c tailwind.config.js -o server/assets/stylesheet.css -w

templ-watch: # Build/rebuild all `templ` templates files in watch mode
	templ generate --watch

go-run: # Run the code without build step in server mode
	go run . -s

### REPL Helper

repl-test: # Test all REPL project module
	go test ./...

repl-build: # Same as make go-build
	go build .

repl-run: # Run repl with-out build step
	go run .

repl-build-run: # Build the code then run executable file
	go build .
	./interingo

repl-clean: # Remove repl build file
	rm -f interingo

### Helper

server-clean: # Remove all frontend build file
	rm -f server/index_templ.go
	rm -f server/assets/stylesheet.css

clean: server-clean repl-clean # Remove all build file

help: # Show this help
	@cat Makefile | \
		grep -E '^[^.:[:space:]]+:|[#]##' | \
		sed -E 's/:[^#]*#([^:]+)$$/: #:\1/' | \
		sed -E 's/([^.:[:space:]]+):([^#]*#(.+))?.*/    make \1\3/' | \
		sed -E 's/[#][#]# *(.+)/# \1/' | \
		column -ts: -L -W2
