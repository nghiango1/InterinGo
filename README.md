> This is a attempt on readding and reproducing the book "Writing an interpreter in Go" step by step. While also trying to use go lang as a language to get famillar with the tool it provided.

# Basic programming concept

A parse and a tokenizer in Interpreter is just a big string parsing

# Personal walkthrough note

## Auto pillot

Chapter 2:
When reading the book, i have almost no context how the interpreter parser acctually work. While this is a subject in my collegue year, what we mostly done is "running" program parsing work on paper (drawing context tree)

Being dump, while dealling with a lot of testing driven coding as how the book appoarch the problem. I ratther copy most of testing code (and ast oop class) provided as is and focus on the acutual main thing program handling, which is parsing the string into each ast class components.

This make thing easier to follow, while making me able to try ahead writing my own code before looking at solution in the book. Most of the code repeat it self quite a lot, but here is what I have done my self

- `main.repl` update atter `program.Parser` and `ast.Node.String()` start working: With `String()`, instead of just print out token literal, I focus on looking how the current code handle each statement, and print out the expected String version (resemble the input). Not thing impressed, just send the `stmt.String()` instead of loop though all and print the `token`. But here are some noteable thing:
    - With `let x = 1` input, repl return just `let x`: We haven't handle the expression yet, so it left a `//TODO` there, I focus on it right after and delay reading the book any further and fail to do so. Which leave the code sending `;` prefix not found error for like 10 commit until I finnaly give up and read the solution
    - With `let x` input, no error is sent out: Just a knowledge check, `fmt.Errorf` not send thing to os.stdout, it need to be followed with `fmt.Println(err.Error())`, read the doc and fix this no time.

- LetParse handing expression after `=` token: It just calling the `expressionParse` function, and assign them into the `ast.LetStatement.Value`. I come and fill this right after finish Prefix and Infix Expression Statement parsing section (before If-Else-parsing, etc..). The thing is that, expression parsing not extending after the end of statement `;` token.  I have a hard time and need to look at what the book provided, which need me to advantage through using `parse.peekNextToken()` check and `parse.NextToken()`. Using git log, you will find me using `parse.currTokenIs()` instead, which falling the job, and lead to program parsing error where is don't understand what to do with `;` token.

- IfElseExpressionParse, BlockStatement ast node, BlockStatementParse: It just repeat it self, now having all the knowledge from previous session, thing get easier for me to write it before all the readding.

- FunctionLiteral: Just tried it, things look quite simmlar but not really, i kept it as is until it break (the last ',' skip in function parameters parse seem weird but it ok)

Here is how REPL output with current code
```
$./main
Hello ylong! This is the Monkey programming language!
Feel free to type in commands
>> fn(x,y) { if (x > y) {return x} else {return y} }
Statement fn(x, y) {if (x > y) {return x} else {return y} ;};
```

