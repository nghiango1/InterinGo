# How to use

Build the program and get `./interingo` file executable or download [Released binary](https://github.com/nghiango1/InterinGo/releases). You can run `./interingo -h` to get help on runner flag directly (for TLDR folks)

## REPL mode

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

## File mode

> This mode have the highest piority, so don't expect server, or REPL running along with `-f` flag

Running `./interingo` executable with `-f` flag.

```sh
./interingo -f <file-location>
```

Unknow what to do yet, use test code in 'test/' directory as your start point. Every file contain comment for expected output in the top to make sure you don't get lost

```sh
./interingo -f test/return-01.iig
```

## Server mode

> As expected, who know what you got if they can't just test it directly on the browser

Running `./interingo` executable with `-s` flag

```sh
./interingo -s
```

You can also specify listen address with `-l` flag or it will default to `0.0.0.0:8080`

```sh
./interingo -s -l 127.0.0.1:4000
```

## Verbose output

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
