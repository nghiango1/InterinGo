> This is a attempt on readding and reproducing the book "Writing an interpreter in Go" step by step. While also trying to use go lang as a language to get famillar with the tool it provided.

# Basic programming concept

A parse and a tokenizer in Interpreter is just a big string parsing

# Personal walkthrough note


## Chapter 2:

### Auto pillot:

When reading the book, i have almost no context how the interpreter parser acctually work. While this is a subject in my collegue year, what we mostly done is "running" program parsing work on paper (drawing context tree)

Being dump, while dealling with a lot of testing driven coding as how the book appoarch the problem. I ratther copy most of testing code (and ast oop class) provided as is and focus on the acutual main thing program handling, which is parsing the string into each ast class components.

This make thing easier to follow, while making me able to try ahead writing my own code before looking at solution in the book. Most of the code repeat it self quite a lot, but here is what I have done my self

### Writing ahead:

I did some change, or rather say I did done some part on my own without looking at the code

- `main.repl` update atter `program.Parser` and `ast.Node.String()` start working: With `String()`, instead of just print out token literal, I focus on looking how the current code handle each statement, and print out the expected String version (resemble the input). Not thing impressed, just send the `stmt.String()` instead of loop though all and print the `token`. But here are some noteable thing:
    - With `let x = 1` input, repl return just `let x`: We haven't handle the expression yet, so it left a `//TODO` there, I focus on it right after and delay reading the book any further and fail to do so. Which leave the code sending `;` prefix not found error for like 10 commits until I finnaly give up and read the solution
    - With `let x` input, no error is sent out: Just a knowledge check, `fmt.Errorf` not send thing to os.stdout, it need to be followed with `fmt.Println(err.Error())`, read the doc and fix this no time.

- LetParse handing expression after `=` token: It just calling the `expressionParse` function, and assign them into the `ast.LetStatement.Value`. I come and fill this right after finish Prefix and Infix Expression Statement parsing section (before If-Else-parsing, etc..). The thing is that, expression parsing not extending after the end of statement `;` token.  I have a hard time and need to look at what the book provided, which need me to advantage through using `parse.peekNextToken()` check and `parse.NextToken()`. Using git log, you will find me using `parse.currTokenIs()` instead, which falling the job, and lead to program parsing error where is don't understand what to do with `;` token.

- IfElseExpressionParse, BlockStatement (ast type and parse), BlockStatementParse: It just repeat it self, now having all the knowledge from previous session, thing get easier for me to write it before all the readding.
    - BlockStatement: I realize it goes to infinity loop when using REPL (we have a nice test case coverage here). This being fix in no time thought

- FunctionLiteral: Just tried it, things look quite simmlar but not really, i kept it as is until it break (the last ',' skip in function parameters parse seem weird but it ok)

- CallExpression: Well, I really unsure about what Call should be, so I at least try to read a litte as possible, while writing code reffering how I understand the parse should do to handle function call. By kept readding, my code desn't work with the precedence, thus the parse can't combine function call with-in an equation expression (The test still pass tho)
    Here is what I have done shown in REPL
    ```
    >> let x = add(x,3) + 5 * (add(2,4) + 4)
    Parser has 3 errors
    Parser error: "no prefix parse function for + found"
    Parser error: "Expected next token to be ), got + instead"
    Parser error: "no prefix parse function for ) found"
    ```

> Final touch at the end of chapter 2. I rewrite some of my wrote REPL code with thing similar from the book just because i thing it a bit more pretty and could save me some time prepare for Chapter 3 ahead
> Still, i do think that, current implement of REPL treat each line as a seperated program isn't a good thing. We will want to keep and re-use our initialized function and variable

## Chapter 3

I don't know much and need to follow the book in most part. Sample REPL output is here

```
(ins)→ ./main
Hello ylong! This is the Monkey programming language!
Feel free to type in commands
>> 55
55
>> 4 == 4
true
>> if (5 * 5 + 10 > 34) { 99 } else { 100 }
99
```

### Auto pillot:

Still, test are meant for some one already know what they need. Looking too deep into test isn't helping that well, I skip and the most of the code.

Eval handling code: While i do try to done them by my self, it is way way different from the book solution and end up getting thing wrong in the most part. It is a normal thing for me to go back and delete all my code then replace by the book solution, thus make this section look identically to the provided code.
- Program `ast` structure: Program -> Statement**s** -> Each statement type way of Eval handling, is not supprising. I still need to look and follow their way implement IntLiteral first and given with 
- If then else: I got wrong with `Eval(ie.Consequense)` -> `Eval(ie.Condition)`, it seem crazy to spot it without any debuger for me i think. Maybe it actually better to get more skill in loging and finding error quickeri.
- Return value: I honestly don't know anything about why it need a struct, and why thing need to be in that exact order. In my best guess, we want to seperating the type of the statement. By casting type to `ReturnValueObject`, we can then know when to stop evaluating `[]Statements` and return the "wrapped" value immediately
    >  We then wrap the result of this call to Eval in our new object.ReturnValue so we can keep track of it.
- Error:
    - The error test with code block object isn't working as intended. I have to change it from 
        ```
        if (10 > 1) {
            if (10 > 1) {
                return true + false;
            }
            return 1;
        }
        ```
        to 
        ```
        if (10 > 1) {
            if (10 > 1) {
                true + false;
                return 10;
            }
            return 1;
        }
        ```
        to generate right error value.
    - We not handle divide by zero error yet, so I implementing it too
        ```go
        // function evalIntergerInfix
        ...
            case "/":
                if rightVal == 0 {
                    return newError("divide by zero: %d %s %d", leftVal, operator, rightVal)
                }
                return &object.Integer{Value: leftVal / rightVal}
        ```
        It also seem that, we can't really pass our Error to AST parent node when Evaluation, making it impossible to handle Nested expression, eg `2/0 > 3`, which let the Expression evaluation kept running. We can try to handle everytime we call Eval without directly return it (or rewrite Eval function to handle it directly, this is not ideal as it lead to even more rewrite), effectively prevent Error object being used as a Value parameter in other evaluation process.
        ```go
        case *ast.PrefixExpression:
            right := Eval(node.Right)
            // Check if Error the then return
            if right != nil && right.Type() == object.ERROR_OBJ {
                return right
            }
            // Else, keep evaluating
            return evalPrefixExpression(node.Operator, right)
        ```
- File mode: Currently, REPL not support multi-line input. Instead of handle that directly, a file input should work just fine for both testing and using the language. I added varius flag for the `main.go` to make it a better CLI program
    - Comment: With file input available, commenting is another concept that needed to be handle
- Server mode: Not related to the book at all, but there is some change I have to made so REPL work well on the web
    - Print: All possible print-out function will using `io.Writer` as interface, no more stdout
- Environment: I use a local packet variable to store `env`, this way it being use in both Server mode and CLI mode
- Function object: This is quite weird, implementing object.Object.Inspect() vs ast.Node.String() seem do the same thing
