# interpreter-in-go

"interprester-in-go" or InterinGo (for short) is a new interpreter language. It can be run in 3 mode
- REPL mode: Which stand for read-evaluation-print-loop, similar to `python`
- File mode: Execute code as input from file
- Server mode: Which have a pretty UI for REPL on a HTTP Server
- Verbose mode: Tell more infomation about Lexer, Parse, Evaluation process

## Why

To challenge my knowledge with `go` language and advanced (interpreter) concept.

## How to use

> Build the program and get `./main` file executable or download Released binary
> Run `./main -h` to get help on runner flag directly (for TLDR folks)

### REPL mode

Running `./main` executable normaly
```sh
./main
```

And you should have been welcome with this
```
→ ./main
Hello <username>! This is the InterinGo programming language!
Type `help()` in commands for common guide
>> 
```

### File mode

> This have the highest piority, so don't expect server, or REPL running 

Running `./main` executable with `-f` flag.
```sh
./main -f <file-location>
```

Unknow what to do yet, use test code in 'test/' directory as your start point. Every file contain comment for expected output in the top to make sure you don't get lost
```sh
./main -f test/return-01.iig
```


### Server mode

> As expected, who know what you got if they can't just test it directly on the browser

Running `./main` executable with `-s` flag
```sh
./main -s
```

You can also specify listen address with `-l` flag or it will default to `0.0.0.0:8080`
```sh
./main -s -l 127.0.0.1:4000
```

### Verbose mode

> Canbe used with any mode

Start with the `-v` flag
```sh
./main -v
```

You can also running `toggleVerbose()`command in InterinGo REPL to enable/disable it on the fly


## The "interprester-in-go" language syntax:

### Keyword

Here is the list of all keyword that already reserved and cannot be used as variable names
```
fn      let      if      else
return  true     false
```

"interpreter-in-go" is a case-sensitive language, meanning that `fn` and `FN` is a different names. making `FN` can be use as a new name to assign a variable

### Statement

Single line statement can end with `;` or not
```
>> let x = 1
>> let y = 2;
```

Double statement in one line can be seperated by a semicolon `;` for ensuring the parse to work correctly.
```
let x = 1; let y = 2;
```

> In-case there is no semicolon, it parse left to right and try to match statement inorder until error occurs 

A block statement (syntax of creating a new function)
```
let f = fn (x,y) {
   return x + y; 
}
```

### Value, Variable, Type initiation and assignment

"interpreter-in-go" is a dynamicly type language (which mean it not have type)

Variable can be initialized and assign with this syntax. 
```
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
```
let x = 1 + 3
let y = x + (-3) * 5 (4 - 9)
let z = NULL
let a = true
let b = !(x >= y)
let c = add(3, 4)
```

Function assign used 
```
let f = add(x,y) { return x + y }
```

### Logic flow control

"interpreter-in-go" support if then else control flow.
```
let x = 1
let y = 2
if ( x >= y ) { let x = 1 } else { let y = 2 }
```


## Build

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

> To make it simplier, the project still contain the build file of templ and tailwindcss

Build the code with

```sh
go build .
```

### Run

Run the program REPL with

```sh
./main
```

### Server mode

InterinGo language can be serve as a http server

> Read more on SERVER.md

```sh
./main -s
```
