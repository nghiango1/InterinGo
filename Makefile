TAILWIND_CLI := tailwindcss-linux-x64

.PHONY: all build build-run tailwind-build tailwind-watch templ-build templ-watch help go-build run

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
	./main -s

### Development helper

tailwind-watch: # Start tailwind in watch mode - looking for change and rebuild, output file server/assets/stylesheet.css
	$(TAILWIND_CLI) -c tailwind.config.js -o server/assets/stylesheet.css -w

templ-watch: # Build/rebuild all `templ` templates files in watch mode
	templ generate --watch

go-run: build run # Run the code without build step in server mode
	go run . -s

### Helper

help: # Show available targets
	@cat Makefile | \
		grep -E '^[^.:[:space:]]+:|[#]##' | \
		sed -E 's/:[^#]*#([^:]+)$$/: #:\1/' | \
		sed -E 's/([^.:[:space:]]+):([^#]*#(.+))?.*/  make \1\3/' | \
		sed -E 's/[#][#]# *(.+)/\n    (\1)\n/' | \
		sed 's/$$/:/' | \
		column -ts:
