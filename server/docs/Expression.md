# Assignment

"interpreter-in-go" is a dynamicly type language (which mean it not have type)

Variable can be initialized and assign a value (expression, variable, function call, ...) to it with this syntax.

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
let y = x + (-3) * 5 (4 - 9)
let z = NULL
let a = true
let b = !(x >= y)
```
