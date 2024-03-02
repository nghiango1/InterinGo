# Return statement

> Control flow is a form of changing how and which statement will be Executed by the Interpreter.

Return statement can stop the Evaluation process imedietly. Let start with a example, run two command bellow simultaneously

```igg
let x = 1;
let y = 2;
return x + y;
let x = 3;
let x = 4;
```


```igg
return x;
```

Here we can see the first command have `return x + y;` statement in beetween a set command. Causing the Interpreter to stop the Evaluation process and return value `3`. Nullify (or skipping) both let function after that, cause `x` value remain `1`.

The next command show the value of `x`, which is `1`. Telling us the Return statement working exactly in expected behavior.

## Syntax

Return statement start with `return` keyword, and following by an Expression statement

> Expression statements: Everything expect a Return statement or a Let statement is a Expression statement

```igg
return 2;
```

## Example

Return in a Function

```igg
let add = fn (x, y) { return x + y; };
return add(5, 4);
```
Return in a complex If expression

```igg
let add = fn (x, y) { return x + y; };

let y = if (6 > 3) {
    let x = 2;
    let y = 2;
    return add(x, y);
};

return y;
```
