# interpreter-in-go

"interprester-in-go" or InterinGo (for short) is a new interpreter language, come with [LSP](lsp-interingo/) and [highlighter](tree-sitter-interingo/) for neovim. It can be run in 3 mode
- REPL mode: Which stand for read-evaluation-print-loop, similar to `python`
- File mode: Execute code as input from file
- Server mode: Which have a pretty UI for REPL on a HTTP Server

## Why

To challenge my knowledge with `go` language and advanced (interpreter) concept. I also set up a http server to public InterinGo interpreter, that you can access evaluating the [language right now](https://nghiango.asia/).

## Techstack:

Front-end
- `templ` for html template
- `tailwind.css` for styling
- `htmx`for server REPL API access
- `javascript` for some dynamic UI rendering

Back-end:
- `golang` standard library for http-request/server handling
- `readline` for CLI

## How to use

> Build the program and get `./interingo` file executable or download Released binary
> Run `./interingo -h` to get help on runner flag directly (for TLDR folks)

### REPL mode

Running `./interingo` executable normaly

```sh
./interingo
```

And you should have been welcome with this

```sh
$ ./interingo
Hello <username>! This is the InterinGo programming language!
Type `help()` in commands for common guide
>>
```

### File mode

> This have the highest piority, so don't expect server, or REPL running

Running `./interingo` executable with `-f` flag.

```sh
./interingo -f <file-location>
```

Unknow what to do yet, use test code in 'test/' directory as your start point. Every file contain comment for expected output in the top to make sure you don't get lost

```sh
./interingo -f test/return-01.iig
```

### Server mode

> As expected, who know what you got if they can't just test it directly on the browser

Running `./interingo` executable with `-s` flag

```sh
./interingo -s
```

You can also specify listen address with `-l` flag or it will default to `0.0.0.0:8080`

```sh
./interingo -s -l 127.0.0.1:4000
```

### Verbose output

Tell more infomation about Lexer, Parse, Evaluation process via REPL output

Start with the `-v` flag

```sh
$ ./interingo -v
```

Or using `toggleVerbose()`command in InterinGo REPL to enable/disable it

```sh
$ ./interingo
>> toggleVerbose()
```

## The "interprester-in-go" language syntax:

You can go to [server/docs/](server/docs/) or read live [docs website](https://nghiango.asia/docs)

## Build - Minimal REPL build

### Prerequisite

Go
- Install go version manager [`gvm`](https://github.com/moovweb/gvm)

    ```sh
    sudo apt-get install bison
    bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
    ```
- `cd` usually not working, so better just delete it.

    ```sh
    vi ~/.gvm/scripts/gvm-default
    # Delete the last line of `gvm-default` file - Which change cd functionality
    ```
- Install latest (currently at v1.22.0) version from binary file and set it as default

    ```sh
    gvm install go1.22.0 -B
    gvm use go1.22.0 --default
    ```

### Build

> To make build process simplier who, the project still contain the build file of templ and tailwindcss. Read more on how to build the front-end in `SERVER.md` file

Build the code with

```sh
go build .
```

### Test

Test all module with

```sh
go test ./...
```

## Build - Full build with Server front-end

All server source file is in `/server/` directory, which need special handle for `templ` files - containing frontend code. This require extra build tool and generating code step. `Makefile` is add to help handle these process

### Prerequire

Install go-lang latest version, currently go 1.22.0

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

Install cmake

```sh
apt-get -y install make
```

### Using `Make` tools

#### Build

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

## Build - LSP, Highlight

LSP and Treesitter highlight can be overkill as it per need specific text editor configuration (some even need specific plugin/extension configuration). Currently, I only setup a Neovim config as Local, as provide them as a fully working component with one click install seem like an overkill. Following each correspond README.md to properly config and use them in neovim if you interested
- Read LSP how to build and use [here](lsp-interingo/README.md) 
- Read treesitter highlight build [here](tree-sitter-interingo/README.md)

I at least can provide some pictures to prove it work:
- In markdown file: [Placeholder]
- In specific `*.iig` file extension: [Placeholder]
- In LSP auto format: [Placeholder]
