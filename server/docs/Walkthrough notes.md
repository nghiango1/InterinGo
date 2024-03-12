
# Personal walkthrough note

A parse and a tokenizer in Interpreter is just a big string parsing

> This is a personal note while i go through the book of How to write an interpreter in Go incase I actually need to tell a story on how I get here

## Lexer and Parser, what it this

### Auto pillot

When reading the book, i have almost no context how the interpreter parser acctually work. While this is a subject in my collegue year, what we mostly done is "running" program parsing work on paper (drawing context tree)

Being dump, while dealling with a lot of testing driven coding as how the book appoarch the problem. I choose to copy most of testing code (and some ast oop class) provided as is and focus on the actual problem handling, which is parsing the string into each ast class components. This make thing easier to follow, while making me able to try ahead writing my own code before looking at solution in the book. 

### Self aware and Writing ahead:

I did some change, or rather say I did done some part on my own without looking at the provided solution

#### `main.repl` - REPL but it is a Read Parse Print Loop now

After `program.Parser` and `ast.Node.String()` start working, with `String()`, instead of just print out token literal, I focus on looking how the current code handle each statement, and print out the expected String version (resemble the input).

Not thing impressed, just send the `stmt.String()` instead of loop though all and print the `token`. But here are some noteable thing:
- With `let x = 1` input, repl return just `let x`: We haven't handle the expression yet, so it left a `//TODO` there, I focus on it right after and delay reading the book any further and fail to do so. Which leave the code sending `;` prefix not found error for like 10 commits until I finnaly give up and read the solution
- With `let x` input, no error is sent out: Just a knowledge check, `fmt.Errorf` not send thing to os.stdout, it need to be followed with `fmt.Println(err.Error())`, read the doc and fix this no time.

> After Eval() start working, I implemented verboseToggle() to REPL to check how my Lexer and Parse work

#### LetStatementParse handing expression

Right after knowing how to handle Prefix and Infix Expression Statement parsing section, I come and fill `//TODO` (wich I just fail),  it just calling the `expressionParse` function, and assign them into the `ast.LetStatement.Value`. How that go?

The thing is that, when calling expression parsing function, `Lexer` is not jump over the end of statement - `;` token. The parser panic and tell me a string that start with `;` isn't match with any defined statement available

I have a hard time and need to look at what the book provided. So I need to advantage through using `parse.peekNextToken()` to check `;` and run `parse.NextToken()` to skip it. Using git log, you will find me using `parse.currTokenIs()` instead, which falling the job and lead to parsing error where REPL don't know what to do with `;` token.

#### IfElseExpressionParse, BlockStatement, BlockStatementParse

It just repeat it self, now having all the knowledge from previous session, thing get easier for me to write it before all the readding.

A while after implemented BlockStatement, I realized it goes to infinity loop when using REPL. Test case coverage seem off here, but it was fix in no time.

#### FunctionLiteral

Just tried it, things look quite simmlar but not really, i kept it as is until it break. To be more percise, I skip the last ',' function parameters parse (or it mean I allow trailling comma in function parameter initilization). It seem weird for a language to have it, but I clearly want it in dictinary or JSON object initilized though.

Now, even my `;` or `\n` (newline) at end of statement isn't a hard force thing, so parse allow something like this:

```iig
let x = 1 let y = 2
x
y
```

#### CallExpression

Well, I really unsure about what Call should be, so I at least try to read a litte as possible, while writing code reffering how I understand the parse should do to handle function call. By kept on readding, my code doesn't work with the precedence, thus the parse can't combine function call with-in an equation expression. The test coverage, again, fail to detect this behavior

Here is what I have done shown in REPL
```iig
let x = add(x,3) + 5 * (add(2,4) + 4)
```

At that time, REPL output thing like this
```sh
> let x = add(x,3) + 5 * (add(2,4) + 4)
Parser has 3 errors
Parser error: "no prefix parse function for + found"
Parser error: "Expected next token to be ), got + instead"
Parser error: "no prefix parse function for ) found"
```

#### Lexer, Parse finish line

Final touch at the end of writing parser and lexer (chapter 2 or the books). I rewrite some of my writen REPL code with solution from the book. This is because that I could save me some time when implementing Evaluation function (Chapter 3 of the book) awaiting ahead.

Still, here is thing that i want to address later:
- Current implement of REPL seem to treat each line as a seperated program isn't a good thing. We will want to keep and re-use our initialized function and variable
    > I later know that this is Wrong
- REPL is the worst CLI experience I could ever have
    > I found out about Readline package, used it, now the final REPL is greate
- No bit operation `>>, <<, ^, &`, no great/less and equal `>=, <=` and some other thing implemented yet. I let this go for now, as It is not a new thing.

## Evaluation parsed program 

### Auto pillot

Again, I don't know much and need to follow the book in most part. Sample REPL output is here Still, test are meant for some one already know what they need. Looking too deep into test isn't helping that well, I skip and the most of the code.

### Self aware and Writing ahead

#### Eval a program - implementation of AST tree structure

I try to done them by my self and end up getting things wrong in the most part with error throwing every where. Okay, let just go back and delete all my code and replace by the book solution.

Program `ast` structure: Program -> Statement**s** -> Each statement type need a different way of Eval handling, (supprising observation?). I still need to look and follow solution implement  

#### Evaluation If then else

Dumb me wrote `Eval(ie.Consequense)` instead of `Eval(ie.Condition)`, causing me ton of debug time.

Also it seem crazy to spot it without any debuger for me i think, Maybe it actually better to get more skill in loging and finding error quicker.

#### Return value

I honestly don't know anything about it. Why it need a struct? And why thing need to be in that exact order?

In my best guess, we want to seperating the type of the statement. By casting type to `ReturnValueObject`, we can then know when to stop evaluating process immediately and return the value.

This way, our new object.ReturnValue type it self is the way to keep track of it.

#### Error:

The error test provided with BlockStatemnt isn't working as intended. I have to change it to check if Evaluation process can stop with Errors inside a BlockStatement properly

```iig
if (10 > 1) {
    if (10 > 1) {
        return true + false;
    }
    return 1;
}
```
to 
```iig
if (10 > 1) {
    if (10 > 1) {
        true + false;
        return 10;
    }
    return 1;
}
```

Also, divide by zero error isn't being handle yet (causing actual error and end REPL program), so I implementing it too
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

I just being ahead of my self here, but this later is mentioned by the book

#### File mode

Currently, REPL not support multi-line input. Instead of handle that directly, a file input should work just fine for both testing and using the language. I added varius flag for the `main.go` to make it a better CLI program

I also implementing Comment. It skip thing in Lexer level (not Parse). With file input available, commenting is must have concept (for me, i like commenting) to be add into the language.

#### Server mode

It just not related to implementing Interpreter it self. There is just nothing to tell me doing this, but who know what you have done without a website throwed at their face.

I do need to change so REPL work well on the web like all print function using `io.Writer` interface

#### Environment

I use a local packet variable to store `env`, this way it being use in both Server mode and CLI mode

#### Function object and Function call evaluation

My current implementation pass global environment into Function it self, which mean there is currently no local scope.
- While this could be change quite easily by introducing another variable in Function Object struct, but for closure to work (or multi nested function local variable), every env function have to be available, which we also need a parrent function in our function object.
- Also, we not have different token to binding GLOBAL/LOCAL variable just yet. As this could be a lot of work, so it just better avoiding these problem for this moment, and only have global variable available in InterinGo language.

Because of this, Arg that passed into a function now directly change the indentifier as a global variable. Causing this code return `15` instead.

```iig
let add = fn(x,y) {
    return x + y;
}
add(5 + 5, add(5, 5))
```

We can do some calcuation here:
- First `add` call set `x = 10`
- Second `add` call set `x = 5` and `y = 5`, overiding `x` value and return 10
- Finally, we have `5 + 10 = 15`

So, we want a seperated arg in function call store from our global env variable. Turn out to change this behavior is just about ~5 line of code, some minor change and all above overhead is just plain wrong, I just need to implement a better env structure.

Also, after keep reading, I known that:
- My `evalCallExpression` implementation, which try to handle `ast.FunctionLiteral.Function` base on two case: is it a `ast.Identifier` or a `ast.FunctionLiteral`? Just use `Eval` and we have `object.Function` in return. Eh, way smaller and nicer.
- Return value can cause chain reaction and stop main Evaluation flow immedietly, which isn't cover by test. So I add this:

    ```iig
    let add = fn(x, y) {
        return x + y
    }
    add(2, 4);
    let x = 2;
    return x;
    ```

    This should return 2

#### High ordered function (closure) support

My implement on eval function call have wrong environment value. That instead update that function object environment, i update my current (what ever scope that the function call was evoke) environment. Causing `ERROR: undefined identifier <variable_name>`. Cost me a while to debug and identify the problem

```diff
-func evalFunctionObject(fo *object.Function, args []ast.Expression, env *object.Environment) object.Object {
+func evalFunctionObject(fo *object.Function, args []ast.Expression) object.Object {
    numOfFuncParam := len(fo.Parameters)
    numOfArgs := len(args)
    if numOfArgs != numOfFuncParam {
        return newError("Function take %d agrument but %d are given", numOfArgs, numOfFuncParam)
    }
 
-   encloseEnv := object.NewEnclosedEnvironment(env)
+   encloseEnv := object.NewEnclosedEnvironment(fo.Env)
    for i := 0; i < numOfFuncParam; i++ {
-       argValue := Eval(args[i], env)
+       argValue := Eval(args[i], fo.Env)
        if isError(argValue) {
            return argValue
        }
```

#### Final line

implementing object.Object.Inspect() vs ast.Node.String() almost done the same thing for object,Function and ast.FunctionLiteral

Nested map structure as environment is too expensive. I do thing back on how C using stack for handle memory after this, could it be replicate in InterinGo?

## Language extending - Tool chain

- Syntax highlighting with Rule or Tree-sitter?
