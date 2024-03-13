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

const
  newline = '\n',
  terminator = choice(newline, ';', '\0');

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
    $.comment_line,
  ],

  rules: {
    // TODO: add the actual grammar rules
    source_file: $ => repeat(
      $._statement
    ),

    comment_line: $ => seq(
      '//',
      /(\\+(.|\r?\n)|[^\\\n])*/
    ),

    _statement: $ => seq(
      choice(
        $.let_statement,
        $.return_statement,
        $.expression_statement
      ),
      optional(terminator)
    ),

    let_statement: $ => seq(
      'let',
      field('name', $.identifier),
      '=',
      field('value', $.expression)
    ),

    return_statement: $ => seq(
      'return',
      $.expression
    ),

    expression_statement: $ => seq(
      $.expression
    ),

    expression: $ => choice(
      $.infix_expression,
      $.prefix_expression
    ),

    prefix_expression: $ => choice(
      prec(precedences_value.PRIORITY, $.if_expression),
      prec(precedences_value.PRIORITY, $.function_literal),
      prec(precedences_value.PRIORITY, seq('(', $.expression, ')')),
      prec(precedences_value.LOWEST, $.identifier),
      prec(precedences_value.LOWEST, $.interger_literal),
      prec(precedences_value.LOWEST, $.boolean),
      prec(precedences_value.PREFIX, seq('-', $.expression)),
      prec(precedences_value.PREFIX, seq('!', $.expression)),
    ),

    infix_expression: $ => {
      const helper = function AddingFieldNameToInfixExpression(l, o, r) {
        return seq(
          field('left', l),
          field('operator', o),
          field('right', r)
        )
      }

      return choice(
        prec.left(precedences.PLUS, helper($.expression, '+', $.expression)),
        prec.left(precedences.MINUS, helper($.expression, '-', $.expression)),
        prec.left(precedences.SLASH, helper($.expression, '/', $.expression)),
        prec.left(precedences.ASTERISK, helper($.expression, '*', $.expression)),
        prec.left(precedences.EQ, helper($.expression, '==', $.expression)),
        prec.left(precedences.NOT_EQ, helper($.expression, '!=', $.expression)),
        prec.left(precedences.GT, helper($.expression, '>', $.expression)),
        prec.left(precedences.LT, helper($.expression, '<', $.expression)),
        prec.left(precedences.LPAREN, $.call_expression)
      )
    },

    identifier: $ => /[a-zA-Z_]+/,

    interger_literal: $ => /\d+/,

    boolean: $ => choice(
      'true',
      'false'
    ),

    if_expression: $ => seq(
      'if',
      '(',
      field('condition', $.expression),
      ')',
      field('consequence', $.code_block),
      optional(
        seq(
          'else',
          field('alternative', $.code_block)
        ),
      )
    ),

    call_expression: $ => prec(precedences_value.CALL, seq(
      field('name', $.identifier),
      field('args', $.argument_list)
    )),

    argument_list: $ => seq(
      '(',
      seq($.expression, repeat(seq(',', $.expression))),
      ')'
    ),

    param_list: $ => seq(
      '(',
      seq($.identifier, repeat(seq(',', $.identifier))),
      ')'
    ),

    code_block: $ => seq(
      '{',
      repeat(
        $._statement,
      ),
      '}'
    ),

    function_literal: $ => seq(
      'fn',
      field('param', $.param_list),
      field('body', $.code_block),
    ),
  }
});

