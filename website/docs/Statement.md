# Statement

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

A block statement: Can appeard in syntax of creating a new function or a if expresion control flow

```iig
// creating a new function
let add = fn(x,y) {
   return x+y;
}

// if then else control flow.
let x = 1;
let y = 2;
if ( x>=y ) { return true; } else { return false; };
```
