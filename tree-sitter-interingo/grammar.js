
// const (
// 	_ int = iota
// 	LOWEST
// 	EQUALS
// 	LESSGREATER
// 	SUM
// 	PRODUCT
// 	PREFIX
// 	CALL
// )
//
// This map the operator token to their piority class
// var precedences = map[token.TokenType]int{
// 	token.EQ:       EQUALS,
// 	token.NOT_EQ:   EQUALS,
// 	token.LT:       LESSGREATER,
// 	token.GT:       LESSGREATER,
// 	token.PLUS:     SUM,
// 	token.MINUS:    SUM,
// 	token.SLASH:    PRODUCT,
// 	token.ASTERISK: PRODUCT,
// 	token.LPAREN:   CALL,
// }

const precedences_value = {
  LOWEST: 0,
  EQUALS: 1,
  LESSGREATER: 2,
  SUM: 3,
  PRODUCT: 4,
  PREFIX: 5,
  CALL: 6,
  PRIORITY: 7
}

const precedences = {
  EQ: precedences_value.EQUALS,
  NOT_EQ: precedences_value.EQUALS,
  LT: precedences_value.LESSGREATER,
  GT: precedences_value.LESSGREATER,
  PLUS: precedences_value.SUM,
  MINUS: precedences_value.SUM,
  SLASH: precedences_value.PRODUCT,
  ASTERISK: precedences_value.PRODUCT,
  LPAREN: precedences_value.CALL,
};

module.exports = grammar({
  name: 'interingo',

  extras: $ => [
    /\s|\\\r?\n/,
    $.commentline,
  ],

  rules: {
    // TODO: add the actual grammar rules
    source_file: $ => repeat(
      choice(
        $._statement,
        $.commentline
      )
    ),

    commentline: $ => seq(
      '//',
      /(\\+(.|\r?\n)|[^\\\n])*/
    ),

    _statement: $ => seq(
      choice(
        $.let_statement,
        $.return_statement,
        $.expression_statement
      ),
      choice('\n', ';', '\0')
    ),

    let_statement: $ => seq(
      'let',
      $.identifier,
      '=',
      $.expression
    ),

    return_statement: $ => seq(
      'return',
      $.expression
    ),

    expression_statement: $ => seq(
      $.expression
    ),

    expression: $ => choice(
      prec(1, $.prefix_expression),
      prec(0, $.infix_expression),
    ),

    prefix_expression: $ => choice(
      prec(precedences_value.LOWEST,$.identifier),
      prec(precedences_value.LOWEST,$.interger_literal),
      prec(precedences_value.LOWEST,$.boolean),
      prec(precedences_value.PREFIX, seq('-', $.expression)),
      prec(precedences_value.PREFIX, seq('!', $.expression)),
      prec(precedences_value.PRIORITY, seq('(', $.expression, ')')),
      prec(precedences_value.PRIORITY, $.if_expression),
      prec(precedences_value.PRIORITY, $.function_literal)
    ),

    infix_expression: $ => (
      prec(precedences.PLUS , seq($.expression,'+',$.expression)),
      prec(precedences.MINUS, seq($.expression,'-',$.expression)),
      prec(precedences.SLASH, seq($.expression,'/',$.expression)),
      prec(precedences.ASTERISK, seq($.expression,'*',$.expression)),
      prec(precedences.EQ, seq($.expression,'==',$.expression)),
      prec(precedences.NOT_EQ, seq($.expression,'!=',$.expression)),
      prec(precedences.GT, seq($.expression,'>',$.expression)),
      prec(precedences.LT, seq($.expression,'<',$.expression)),
      prec(precedences.LPAREN, $.call_expression)
    ),

    identifier: $ => /[a-zA-Z]+/,

    interger_literal: $ => /\d+/,

    boolean: $ => choice(
      'true',
      'false'
    ),

    if_expression: $ => seq(
      'if',
      '(',
      $.expression,
      ')',
      $.code_block,
      choice(
        seq(
          'else',
          $.code_block
        ),
        '\0'
      )
    ),

    call_expression: $ => seq(
      $.identifier,
      $.argument_list,
    ),

    argument_list: $ => seq(
      '(',
      repeat(
        $.expression,
      ),
      ')'
    ),

    param_list: $ => seq(
      '(',
      repeat(
        $.identifier,
      ),
      ')'
    ),

    code_block: $ => seq(
      '{',
      repeat(
        choice(
          $._statement,
          $.commentline
        )
      ),
      '}'
    ),

    function_literal: $ => seq(
      'fn',
      $.param_list,
      $.code_block,
    ),
  }
});

