# interpreter-in-go

"interprester-in-go" or InterinGo (for short) is a new interpreter language. It can be run in 3 mode
- REPL mode: Which stand for read-evaluation-print-loop, similar to `python`
- File mode: Execute code as input from file
- Server mode: Which have a pretty UI for REPL on a HTTP Server

## Why

To challenge my knowledge with `go` language and advanced (interpreter) concept. I also set up a http server to public InterinGo interpreter. You can access evaluating the language right now via link in the repo [Github description](https://github.com/nghiango1/hello).

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

### Keyword

Here is the list of all keyword that already reserved and cannot be used as variable names

```iig
fn      let      if      else
return  true     false
```

"interpreter-in-go" is a case-sensitive language, meanning that `fn` and `FN` is a different names. making `FN` can be use as a new name to assign a variable

### Statement

Single line statement can end with `;` or not

```iig
let x = 1
let y = 2;
```

Double statement in one line can be seperated by a semicolon `;` for ensuring the parse to work correctly.

```iig
let x = 1; let y = 2;
```

> In-case there is no semicolon, it parse left to right and try to match statement inorder until error occurs

A block statement (syntax of creating a new function)

```iig
let f = fn (x,y) {
   return x + y;
}
```

### Value, Variable, Type initiation and assignment

"interpreter-in-go" is a dynamicly type language (which mean it not have type)

Variable can be initialized and assign with this syntax.

```iig
let x = 1
```

List of right value expression support:
- Boolean (`true`, `false`)
- Integer (`0`, `1`, ...)
- Identifiers (any string that not a part of the language keyword)
- `NULL`
- Numeric Operation (`+`,`-`, `*`, `/`)
- Comparation (`>`, `<`, `<=`, `>=`, `==`)
- Prefix operation (`!true`, `-(4+3)`)
- Function call(`add(3,4)`)

More example:

```iig
let x = 1 + 3
```
```iig
let y = x + (-3) * 5 (4 - 9)
```
```iig
let z = NULL
```
```iig
let a = true
```
```iig
let b = !(x >= y)
```
```iig
let c = add(3, 4)
```

Function assign used

```iig
let f = add(x,y) { return x + y }
```

### Logic flow control

"interpreter-in-go" support if then else control flow.

```iig
let x = 1; let y = 2; if ( x >= y ) { return true; } else { return false; };
```

## Build - REPL

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

## Build - Server front-end

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
