All InterinGo source file can be use for testing. Which end with-in `.iig` file extension

# Correctly stop evaluation

## Error handling `error-*.iig`

A purposedly error code that use to the how InterinGo REPL handle and showing errors for code and effectively stop executation right after an Error
- `error-01.iig`: First error happen in a if condition expression, where we done a devided by Zero `2/0`
- `error-02.iig`: First error happen in a nested if condition, where we done a devided by Zero `2/0`


## Return statement `return-*.igg`

A return statement mean we can stop the Evaluation and directly return the value as result
- `return-01.igg`: There is a return statement so the rest of the statement isn't being Evaluated any more

## Function

A function return a `return` value can by tricky. We may stop main flow of our Evaluation
- `function-01.igg`: There is a function which return a `return` value, we try to get over that function call and do the rest of the code instruction
