; highlights

[
  "let"
  "return"
  "fn"
  "if"
  "else"
] @keyword

[
  "true"
  "false"
] @constant.builtin

[
  "="
  "-"
  "+"
  "*"
  "/"
  "!="
  "=="
  "<"
  ">"
  "!"
] @operator

(interger_literal) @number

(let_statement
  name: (identifier) @function
  value: (expression
   (prefix_expression
     (function_literal)
   )
  )
)

(call_expression
  name: (identifier) @function
)

(comment_line) @comment
