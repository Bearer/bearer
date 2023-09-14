module.exports = grammar({
  name: 'config_variables',

  extras: $ => [],

  rules: {
    string: $ => repeat(choice(
      $._expansion,
      $.literal
    )),

    literal: $ => prec.left(-2, repeat1(choice(/[^${]+/, /./))),

    variable: $ => prec.left(1, /[^\s${}'"]+/),
    unknown: $ => prec.left(-1, /[^\s{}]([^{}]*[^\s{}])?/),
    _expression: $ => choice($.variable, $.unknown),

    _expansion: $ => choice(
      $._mustache_expansion,
      $._github_actions_expansion,
      $._simple_environment_expansion,
      $._bracketed_environment_expansion
    ),

    _mustache_expansion: $ => seq('{{', optional(/\s+/), $._expression, optional(/\s+/), '}}'),
    _github_actions_expansion: $ => seq('$', $._mustache_expansion),
    _bracketed_environment_expansion: $ => seq('$', '{', optional(/\s+/), $._expression, optional(/\s+/), '}'),
    _simple_environment_expansion: $ => seq('$', $.variable)
  }
});
