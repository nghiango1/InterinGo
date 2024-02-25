# `InterinGo` language - Server mode

Extending funtionality of the project, where I setup a HTTP Server.

Techstack: 
- `templ` for html template
- `tailwind.css` for styling
- `golang` standard library for http-request/server handling

## Why?

It better to have the server public InterinGo interpreter so that others can easier access, evaluating the language i have build.

## How to use

Start the server with default listening on all network interface port 8080 - "0.0.0.0:8080"

```
go run . -s
```

or build and run the executable

```
go build .
./main -s
```

## Build the code

All server source file is in `/server/` directory, which need special handle for `templ` files - containing frontend code. This require extra build tool and generating code step. `Makefile` is add to help handle these process

### Prerequire

Install go-lang latest version, currently go 1.22.0 (my go.mod said 1.21.4 so that work too). I use `gvm` tools for install and manage go version, check `/hello/go/install/README.md` if you need help installing `go`
```sh
gvm install go1.22.0 -B
gvm use go1.22.0 -default
```

Install `templ` tools, learn more in [templ.guide](https://templ.guide/)
```sh
go install github.com/a-h/templ/cmd/templ@latest
```

Download latest `tailwind` CLI standalone tool from their [github](https://github.com/tailwindlabs/tailwindcss/releases/) and put it in to `PATH`. This should be add in `.profile` file
```sh
wget https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.1/tailwindcss-linux-x64
cp tailwindcss-linux-x64 ~/.local/bin
export PATH="$HOME/.local/bin:$PATH"
```

### Using `Make` tools

#### Build mode

I setup Makefile to handle CLI operation, use `make build-run` to rebuild and start the server
- `make` or `make all` or `make help`: Show all option command
- `make build`: Build/Rebuild `templ` template files and generating tailwindcss stylesheet
- `make run`: Run built the server
- `make build-run`: Do both

Example
```sh
make
```

#### Dev mode

Golang doesn't have watch mode, but `templ` and `tailwindcss` have it
- `make tailwind-watch`: Tailwind watch mode - Auto rebuild when files change
- `make templ-watch`: Templ watch mode - Auto rebuild when files change
- `go run . -s` or `make go-run`: Run the server without build
