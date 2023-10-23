#include <tree_sitter/parser.h>

#if defined(__GNUC__) || defined(__clang__)
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wmissing-field-initializers"
#endif

#define LANGUAGE_VERSION 14
#define STATE_COUNT 33
#define LARGE_STATE_COUNT 4
#define SYMBOL_COUNT 23
#define ALIAS_COUNT 0
#define TOKEN_COUNT 11
#define EXTERNAL_TOKEN_COUNT 0
#define FIELD_COUNT 0
#define MAX_ALIAS_SEQUENCE_LENGTH 6
#define PRODUCTION_ID_COUNT 1

enum {
  aux_sym_literal_token1 = 1,
  aux_sym_literal_token2 = 2,
  aux_sym_variable_token1 = 3,
  aux_sym_unknown_token1 = 4,
  anon_sym_LBRACE_LBRACE = 5,
  aux_sym__mustache_expansion_token1 = 6,
  anon_sym_RBRACE_RBRACE = 7,
  anon_sym_DOLLAR = 8,
  anon_sym_LBRACE = 9,
  anon_sym_RBRACE = 10,
  sym_string = 11,
  sym_literal = 12,
  sym_variable = 13,
  sym_unknown = 14,
  sym__expression = 15,
  sym__expansion = 16,
  sym__mustache_expansion = 17,
  sym__github_actions_expansion = 18,
  sym__bracketed_environment_expansion = 19,
  sym__simple_environment_expansion = 20,
  aux_sym_string_repeat1 = 21,
  aux_sym_literal_repeat1 = 22,
};

static const char * const ts_symbol_names[] = {
  [ts_builtin_sym_end] = "end",
  [aux_sym_literal_token1] = "literal_token1",
  [aux_sym_literal_token2] = "literal_token2",
  [aux_sym_variable_token1] = "variable_token1",
  [aux_sym_unknown_token1] = "unknown_token1",
  [anon_sym_LBRACE_LBRACE] = "{{",
  [aux_sym__mustache_expansion_token1] = "_mustache_expansion_token1",
  [anon_sym_RBRACE_RBRACE] = "}}",
  [anon_sym_DOLLAR] = "$",
  [anon_sym_LBRACE] = "{",
  [anon_sym_RBRACE] = "}",
  [sym_string] = "string",
  [sym_literal] = "literal",
  [sym_variable] = "variable",
  [sym_unknown] = "unknown",
  [sym__expression] = "_expression",
  [sym__expansion] = "_expansion",
  [sym__mustache_expansion] = "_mustache_expansion",
  [sym__github_actions_expansion] = "_github_actions_expansion",
  [sym__bracketed_environment_expansion] = "_bracketed_environment_expansion",
  [sym__simple_environment_expansion] = "_simple_environment_expansion",
  [aux_sym_string_repeat1] = "string_repeat1",
  [aux_sym_literal_repeat1] = "literal_repeat1",
};

static const TSSymbol ts_symbol_map[] = {
  [ts_builtin_sym_end] = ts_builtin_sym_end,
  [aux_sym_literal_token1] = aux_sym_literal_token1,
  [aux_sym_literal_token2] = aux_sym_literal_token2,
  [aux_sym_variable_token1] = aux_sym_variable_token1,
  [aux_sym_unknown_token1] = aux_sym_unknown_token1,
  [anon_sym_LBRACE_LBRACE] = anon_sym_LBRACE_LBRACE,
  [aux_sym__mustache_expansion_token1] = aux_sym__mustache_expansion_token1,
  [anon_sym_RBRACE_RBRACE] = anon_sym_RBRACE_RBRACE,
  [anon_sym_DOLLAR] = anon_sym_DOLLAR,
  [anon_sym_LBRACE] = anon_sym_LBRACE,
  [anon_sym_RBRACE] = anon_sym_RBRACE,
  [sym_string] = sym_string,
  [sym_literal] = sym_literal,
  [sym_variable] = sym_variable,
  [sym_unknown] = sym_unknown,
  [sym__expression] = sym__expression,
  [sym__expansion] = sym__expansion,
  [sym__mustache_expansion] = sym__mustache_expansion,
  [sym__github_actions_expansion] = sym__github_actions_expansion,
  [sym__bracketed_environment_expansion] = sym__bracketed_environment_expansion,
  [sym__simple_environment_expansion] = sym__simple_environment_expansion,
  [aux_sym_string_repeat1] = aux_sym_string_repeat1,
  [aux_sym_literal_repeat1] = aux_sym_literal_repeat1,
};

static const TSSymbolMetadata ts_symbol_metadata[] = {
  [ts_builtin_sym_end] = {
    .visible = false,
    .named = true,
  },
  [aux_sym_literal_token1] = {
    .visible = false,
    .named = false,
  },
  [aux_sym_literal_token2] = {
    .visible = false,
    .named = false,
  },
  [aux_sym_variable_token1] = {
    .visible = false,
    .named = false,
  },
  [aux_sym_unknown_token1] = {
    .visible = false,
    .named = false,
  },
  [anon_sym_LBRACE_LBRACE] = {
    .visible = true,
    .named = false,
  },
  [aux_sym__mustache_expansion_token1] = {
    .visible = false,
    .named = false,
  },
  [anon_sym_RBRACE_RBRACE] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_DOLLAR] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_LBRACE] = {
    .visible = true,
    .named = false,
  },
  [anon_sym_RBRACE] = {
    .visible = true,
    .named = false,
  },
  [sym_string] = {
    .visible = true,
    .named = true,
  },
  [sym_literal] = {
    .visible = true,
    .named = true,
  },
  [sym_variable] = {
    .visible = true,
    .named = true,
  },
  [sym_unknown] = {
    .visible = true,
    .named = true,
  },
  [sym__expression] = {
    .visible = false,
    .named = true,
  },
  [sym__expansion] = {
    .visible = false,
    .named = true,
  },
  [sym__mustache_expansion] = {
    .visible = false,
    .named = true,
  },
  [sym__github_actions_expansion] = {
    .visible = false,
    .named = true,
  },
  [sym__bracketed_environment_expansion] = {
    .visible = false,
    .named = true,
  },
  [sym__simple_environment_expansion] = {
    .visible = false,
    .named = true,
  },
  [aux_sym_string_repeat1] = {
    .visible = false,
    .named = false,
  },
  [aux_sym_literal_repeat1] = {
    .visible = false,
    .named = false,
  },
};

static const TSSymbol ts_alias_sequences[PRODUCTION_ID_COUNT][MAX_ALIAS_SEQUENCE_LENGTH] = {
  [0] = {0},
};

static const uint16_t ts_non_terminal_alias_map[] = {
  0,
};

static const TSStateId ts_primary_state_ids[STATE_COUNT] = {
  [0] = 0,
  [1] = 1,
  [2] = 2,
  [3] = 3,
  [4] = 4,
  [5] = 5,
  [6] = 6,
  [7] = 7,
  [8] = 8,
  [9] = 9,
  [10] = 10,
  [11] = 11,
  [12] = 12,
  [13] = 13,
  [14] = 14,
  [15] = 15,
  [16] = 16,
  [17] = 17,
  [18] = 18,
  [19] = 19,
  [20] = 20,
  [21] = 21,
  [22] = 22,
  [23] = 23,
  [24] = 24,
  [25] = 10,
  [26] = 24,
  [27] = 10,
  [28] = 28,
  [29] = 29,
  [30] = 30,
  [31] = 31,
  [32] = 32,
};

static inline bool aux_sym_variable_token1_character_set_1(int32_t c) {
  return (c < '"'
    ? (c < '\r'
      ? (c < '\t'
        ? c == 0
        : c <= '\n')
      : (c <= '\r' || c == ' '))
    : (c <= '"' || (c < '{'
      ? (c < '\''
        ? c == '$'
        : c <= '\'')
      : (c <= '{' || c == '}'))));
}

static bool ts_lex(TSLexer *lexer, TSStateId state) {
  START_LEXER();
  eof = lexer->eof(lexer);
  switch (state) {
    case 0:
      if (eof) ADVANCE(6);
      if (lookahead == '$') ADVANCE(16);
      if (lookahead == '{') ADVANCE(17);
      if (lookahead == '}') ADVANCE(18);
      if (lookahead != 0 &&
          lookahead != '\n') ADVANCE(8);
      END_STATE();
    case 1:
      if (lookahead == '{') ADVANCE(17);
      if (lookahead == '}') ADVANCE(18);
      if (lookahead == '\t' ||
          lookahead == '\n' ||
          lookahead == '\r' ||
          lookahead == ' ') ADVANCE(14);
      if (lookahead != 0 &&
          lookahead != '"' &&
          lookahead != '$' &&
          lookahead != '\'') ADVANCE(11);
      END_STATE();
    case 2:
      if (lookahead == '}') ADVANCE(3);
      if (lookahead == '"' ||
          lookahead == '$' ||
          lookahead == '\'') ADVANCE(12);
      if (lookahead == '\t' ||
          lookahead == '\n' ||
          lookahead == '\r' ||
          lookahead == ' ') ADVANCE(14);
      if (lookahead != 0 &&
          lookahead != '{') ADVANCE(10);
      END_STATE();
    case 3:
      if (lookahead == '}') ADVANCE(15);
      END_STATE();
    case 4:
      if (lookahead == '\t' ||
          lookahead == '\n' ||
          lookahead == '\r' ||
          lookahead == ' ') ADVANCE(4);
      if (lookahead != 0 &&
          lookahead != '{' &&
          lookahead != '}') ADVANCE(12);
      END_STATE();
    case 5:
      if (eof) ADVANCE(6);
      if (lookahead == '\n') ADVANCE(7);
      if (lookahead == '$') ADVANCE(16);
      if (lookahead == '{') ADVANCE(9);
      if (lookahead != 0) ADVANCE(7);
      END_STATE();
    case 6:
      ACCEPT_TOKEN(ts_builtin_sym_end);
      END_STATE();
    case 7:
      ACCEPT_TOKEN(aux_sym_literal_token1);
      if (lookahead != 0 &&
          lookahead != '$' &&
          lookahead != '{') ADVANCE(7);
      END_STATE();
    case 8:
      ACCEPT_TOKEN(aux_sym_literal_token2);
      END_STATE();
    case 9:
      ACCEPT_TOKEN(aux_sym_literal_token2);
      if (lookahead == '{') ADVANCE(13);
      END_STATE();
    case 10:
      ACCEPT_TOKEN(aux_sym_variable_token1);
      if (lookahead == '"' ||
          lookahead == '$' ||
          lookahead == '\'') ADVANCE(12);
      if (lookahead == '\t' ||
          lookahead == '\n' ||
          lookahead == '\r' ||
          lookahead == ' ') ADVANCE(4);
      if (lookahead != 0 &&
          lookahead != '{' &&
          lookahead != '}') ADVANCE(10);
      END_STATE();
    case 11:
      ACCEPT_TOKEN(aux_sym_variable_token1);
      if (!aux_sym_variable_token1_character_set_1(lookahead)) ADVANCE(11);
      END_STATE();
    case 12:
      ACCEPT_TOKEN(aux_sym_unknown_token1);
      if (lookahead == '\t' ||
          lookahead == '\n' ||
          lookahead == '\r' ||
          lookahead == ' ') ADVANCE(4);
      if (lookahead != 0 &&
          lookahead != '{' &&
          lookahead != '}') ADVANCE(12);
      END_STATE();
    case 13:
      ACCEPT_TOKEN(anon_sym_LBRACE_LBRACE);
      END_STATE();
    case 14:
      ACCEPT_TOKEN(aux_sym__mustache_expansion_token1);
      if (lookahead == '\t' ||
          lookahead == '\n' ||
          lookahead == '\r' ||
          lookahead == ' ') ADVANCE(14);
      END_STATE();
    case 15:
      ACCEPT_TOKEN(anon_sym_RBRACE_RBRACE);
      END_STATE();
    case 16:
      ACCEPT_TOKEN(anon_sym_DOLLAR);
      END_STATE();
    case 17:
      ACCEPT_TOKEN(anon_sym_LBRACE);
      if (lookahead == '{') ADVANCE(13);
      END_STATE();
    case 18:
      ACCEPT_TOKEN(anon_sym_RBRACE);
      END_STATE();
    default:
      return false;
  }
}

static const TSLexMode ts_lex_modes[STATE_COUNT] = {
  [0] = {.lex_state = 0},
  [1] = {.lex_state = 5},
  [2] = {.lex_state = 5},
  [3] = {.lex_state = 5},
  [4] = {.lex_state = 2},
  [5] = {.lex_state = 5},
  [6] = {.lex_state = 2},
  [7] = {.lex_state = 5},
  [8] = {.lex_state = 5},
  [9] = {.lex_state = 1},
  [10] = {.lex_state = 5},
  [11] = {.lex_state = 5},
  [12] = {.lex_state = 5},
  [13] = {.lex_state = 2},
  [14] = {.lex_state = 5},
  [15] = {.lex_state = 5},
  [16] = {.lex_state = 5},
  [17] = {.lex_state = 5},
  [18] = {.lex_state = 2},
  [19] = {.lex_state = 5},
  [20] = {.lex_state = 2},
  [21] = {.lex_state = 1},
  [22] = {.lex_state = 1},
  [23] = {.lex_state = 2},
  [24] = {.lex_state = 2},
  [25] = {.lex_state = 2},
  [26] = {.lex_state = 1},
  [27] = {.lex_state = 1},
  [28] = {.lex_state = 2},
  [29] = {.lex_state = 2},
  [30] = {.lex_state = 1},
  [31] = {.lex_state = 0},
  [32] = {.lex_state = 1},
};

static const uint16_t ts_parse_table[LARGE_STATE_COUNT][SYMBOL_COUNT] = {
  [0] = {
    [ts_builtin_sym_end] = ACTIONS(1),
    [aux_sym_literal_token2] = ACTIONS(1),
    [anon_sym_LBRACE_LBRACE] = ACTIONS(1),
    [anon_sym_DOLLAR] = ACTIONS(1),
    [anon_sym_LBRACE] = ACTIONS(1),
    [anon_sym_RBRACE] = ACTIONS(1),
  },
  [1] = {
    [sym_string] = STATE(31),
    [sym_literal] = STATE(2),
    [sym__expansion] = STATE(2),
    [sym__mustache_expansion] = STATE(2),
    [sym__github_actions_expansion] = STATE(2),
    [sym__bracketed_environment_expansion] = STATE(2),
    [sym__simple_environment_expansion] = STATE(2),
    [aux_sym_string_repeat1] = STATE(2),
    [aux_sym_literal_repeat1] = STATE(5),
    [ts_builtin_sym_end] = ACTIONS(3),
    [aux_sym_literal_token1] = ACTIONS(5),
    [aux_sym_literal_token2] = ACTIONS(7),
    [anon_sym_LBRACE_LBRACE] = ACTIONS(9),
    [anon_sym_DOLLAR] = ACTIONS(11),
  },
  [2] = {
    [sym_literal] = STATE(3),
    [sym__expansion] = STATE(3),
    [sym__mustache_expansion] = STATE(3),
    [sym__github_actions_expansion] = STATE(3),
    [sym__bracketed_environment_expansion] = STATE(3),
    [sym__simple_environment_expansion] = STATE(3),
    [aux_sym_string_repeat1] = STATE(3),
    [aux_sym_literal_repeat1] = STATE(5),
    [ts_builtin_sym_end] = ACTIONS(13),
    [aux_sym_literal_token1] = ACTIONS(5),
    [aux_sym_literal_token2] = ACTIONS(7),
    [anon_sym_LBRACE_LBRACE] = ACTIONS(9),
    [anon_sym_DOLLAR] = ACTIONS(11),
  },
  [3] = {
    [sym_literal] = STATE(3),
    [sym__expansion] = STATE(3),
    [sym__mustache_expansion] = STATE(3),
    [sym__github_actions_expansion] = STATE(3),
    [sym__bracketed_environment_expansion] = STATE(3),
    [sym__simple_environment_expansion] = STATE(3),
    [aux_sym_string_repeat1] = STATE(3),
    [aux_sym_literal_repeat1] = STATE(5),
    [ts_builtin_sym_end] = ACTIONS(15),
    [aux_sym_literal_token1] = ACTIONS(17),
    [aux_sym_literal_token2] = ACTIONS(20),
    [anon_sym_LBRACE_LBRACE] = ACTIONS(23),
    [anon_sym_DOLLAR] = ACTIONS(26),
  },
};

static const uint16_t ts_small_parse_table[] = {
  [0] = 4,
    ACTIONS(29), 1,
      aux_sym_variable_token1,
    ACTIONS(31), 1,
      aux_sym_unknown_token1,
    ACTIONS(33), 1,
      aux_sym__mustache_expansion_token1,
    STATE(23), 3,
      sym_variable,
      sym_unknown,
      sym__expression,
  [15] = 4,
    ACTIONS(37), 1,
      aux_sym_literal_token1,
    ACTIONS(39), 1,
      aux_sym_literal_token2,
    STATE(7), 1,
      aux_sym_literal_repeat1,
    ACTIONS(35), 3,
      ts_builtin_sym_end,
      anon_sym_LBRACE_LBRACE,
      anon_sym_DOLLAR,
  [30] = 4,
    ACTIONS(41), 1,
      aux_sym_variable_token1,
    ACTIONS(43), 1,
      aux_sym_unknown_token1,
    ACTIONS(45), 1,
      aux_sym__mustache_expansion_token1,
    STATE(22), 3,
      sym_variable,
      sym_unknown,
      sym__expression,
  [45] = 4,
    ACTIONS(49), 1,
      aux_sym_literal_token1,
    ACTIONS(52), 1,
      aux_sym_literal_token2,
    STATE(7), 1,
      aux_sym_literal_repeat1,
    ACTIONS(47), 3,
      ts_builtin_sym_end,
      anon_sym_LBRACE_LBRACE,
      anon_sym_DOLLAR,
  [60] = 2,
    ACTIONS(57), 1,
      aux_sym_literal_token2,
    ACTIONS(55), 4,
      ts_builtin_sym_end,
      aux_sym_literal_token1,
      anon_sym_LBRACE_LBRACE,
      anon_sym_DOLLAR,
  [70] = 5,
    ACTIONS(9), 1,
      anon_sym_LBRACE_LBRACE,
    ACTIONS(59), 1,
      aux_sym_variable_token1,
    ACTIONS(61), 1,
      anon_sym_LBRACE,
    STATE(15), 1,
      sym_variable,
    STATE(16), 1,
      sym__mustache_expansion,
  [86] = 2,
    ACTIONS(65), 1,
      aux_sym_literal_token2,
    ACTIONS(63), 4,
      ts_builtin_sym_end,
      aux_sym_literal_token1,
      anon_sym_LBRACE_LBRACE,
      anon_sym_DOLLAR,
  [96] = 2,
    ACTIONS(69), 1,
      aux_sym_literal_token2,
    ACTIONS(67), 4,
      ts_builtin_sym_end,
      aux_sym_literal_token1,
      anon_sym_LBRACE_LBRACE,
      anon_sym_DOLLAR,
  [106] = 2,
    ACTIONS(73), 1,
      aux_sym_literal_token2,
    ACTIONS(71), 4,
      ts_builtin_sym_end,
      aux_sym_literal_token1,
      anon_sym_LBRACE_LBRACE,
      anon_sym_DOLLAR,
  [116] = 3,
    ACTIONS(29), 1,
      aux_sym_variable_token1,
    ACTIONS(31), 1,
      aux_sym_unknown_token1,
    STATE(20), 3,
      sym_variable,
      sym_unknown,
      sym__expression,
  [128] = 2,
    ACTIONS(77), 1,
      aux_sym_literal_token2,
    ACTIONS(75), 4,
      ts_builtin_sym_end,
      aux_sym_literal_token1,
      anon_sym_LBRACE_LBRACE,
      anon_sym_DOLLAR,
  [138] = 2,
    ACTIONS(81), 1,
      aux_sym_literal_token2,
    ACTIONS(79), 4,
      ts_builtin_sym_end,
      aux_sym_literal_token1,
      anon_sym_LBRACE_LBRACE,
      anon_sym_DOLLAR,
  [148] = 2,
    ACTIONS(85), 1,
      aux_sym_literal_token2,
    ACTIONS(83), 4,
      ts_builtin_sym_end,
      aux_sym_literal_token1,
      anon_sym_LBRACE_LBRACE,
      anon_sym_DOLLAR,
  [158] = 2,
    ACTIONS(89), 1,
      aux_sym_literal_token2,
    ACTIONS(87), 4,
      ts_builtin_sym_end,
      aux_sym_literal_token1,
      anon_sym_LBRACE_LBRACE,
      anon_sym_DOLLAR,
  [168] = 3,
    ACTIONS(41), 1,
      aux_sym_variable_token1,
    ACTIONS(43), 1,
      aux_sym_unknown_token1,
    STATE(21), 3,
      sym_variable,
      sym_unknown,
      sym__expression,
  [180] = 2,
    ACTIONS(93), 1,
      aux_sym_literal_token2,
    ACTIONS(91), 4,
      ts_builtin_sym_end,
      aux_sym_literal_token1,
      anon_sym_LBRACE_LBRACE,
      anon_sym_DOLLAR,
  [190] = 2,
    ACTIONS(95), 1,
      aux_sym__mustache_expansion_token1,
    ACTIONS(97), 1,
      anon_sym_RBRACE_RBRACE,
  [197] = 2,
    ACTIONS(99), 1,
      aux_sym__mustache_expansion_token1,
    ACTIONS(101), 1,
      anon_sym_RBRACE,
  [204] = 2,
    ACTIONS(103), 1,
      aux_sym__mustache_expansion_token1,
    ACTIONS(105), 1,
      anon_sym_RBRACE,
  [211] = 2,
    ACTIONS(107), 1,
      aux_sym__mustache_expansion_token1,
    ACTIONS(109), 1,
      anon_sym_RBRACE_RBRACE,
  [218] = 1,
    ACTIONS(111), 2,
      aux_sym__mustache_expansion_token1,
      anon_sym_RBRACE_RBRACE,
  [223] = 1,
    ACTIONS(63), 2,
      aux_sym__mustache_expansion_token1,
      anon_sym_RBRACE_RBRACE,
  [228] = 1,
    ACTIONS(111), 2,
      aux_sym__mustache_expansion_token1,
      anon_sym_RBRACE,
  [233] = 1,
    ACTIONS(63), 2,
      aux_sym__mustache_expansion_token1,
      anon_sym_RBRACE,
  [238] = 1,
    ACTIONS(113), 1,
      anon_sym_RBRACE_RBRACE,
  [242] = 1,
    ACTIONS(97), 1,
      anon_sym_RBRACE_RBRACE,
  [246] = 1,
    ACTIONS(101), 1,
      anon_sym_RBRACE,
  [250] = 1,
    ACTIONS(115), 1,
      ts_builtin_sym_end,
  [254] = 1,
    ACTIONS(117), 1,
      anon_sym_RBRACE,
};

static const uint32_t ts_small_parse_table_map[] = {
  [SMALL_STATE(4)] = 0,
  [SMALL_STATE(5)] = 15,
  [SMALL_STATE(6)] = 30,
  [SMALL_STATE(7)] = 45,
  [SMALL_STATE(8)] = 60,
  [SMALL_STATE(9)] = 70,
  [SMALL_STATE(10)] = 86,
  [SMALL_STATE(11)] = 96,
  [SMALL_STATE(12)] = 106,
  [SMALL_STATE(13)] = 116,
  [SMALL_STATE(14)] = 128,
  [SMALL_STATE(15)] = 138,
  [SMALL_STATE(16)] = 148,
  [SMALL_STATE(17)] = 158,
  [SMALL_STATE(18)] = 168,
  [SMALL_STATE(19)] = 180,
  [SMALL_STATE(20)] = 190,
  [SMALL_STATE(21)] = 197,
  [SMALL_STATE(22)] = 204,
  [SMALL_STATE(23)] = 211,
  [SMALL_STATE(24)] = 218,
  [SMALL_STATE(25)] = 223,
  [SMALL_STATE(26)] = 228,
  [SMALL_STATE(27)] = 233,
  [SMALL_STATE(28)] = 238,
  [SMALL_STATE(29)] = 242,
  [SMALL_STATE(30)] = 246,
  [SMALL_STATE(31)] = 250,
  [SMALL_STATE(32)] = 254,
};

static const TSParseActionEntry ts_parse_actions[] = {
  [0] = {.entry = {.count = 0, .reusable = false}},
  [1] = {.entry = {.count = 1, .reusable = false}}, RECOVER(),
  [3] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_string, 0),
  [5] = {.entry = {.count = 1, .reusable = true}}, SHIFT(5),
  [7] = {.entry = {.count = 1, .reusable = false}}, SHIFT(5),
  [9] = {.entry = {.count = 1, .reusable = true}}, SHIFT(4),
  [11] = {.entry = {.count = 1, .reusable = true}}, SHIFT(9),
  [13] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_string, 1),
  [15] = {.entry = {.count = 1, .reusable = true}}, REDUCE(aux_sym_string_repeat1, 2),
  [17] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_string_repeat1, 2), SHIFT_REPEAT(5),
  [20] = {.entry = {.count = 2, .reusable = false}}, REDUCE(aux_sym_string_repeat1, 2), SHIFT_REPEAT(5),
  [23] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_string_repeat1, 2), SHIFT_REPEAT(4),
  [26] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_string_repeat1, 2), SHIFT_REPEAT(9),
  [29] = {.entry = {.count = 1, .reusable = false}}, SHIFT(25),
  [31] = {.entry = {.count = 1, .reusable = false}}, SHIFT(24),
  [33] = {.entry = {.count = 1, .reusable = true}}, SHIFT(13),
  [35] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_literal, 1),
  [37] = {.entry = {.count = 1, .reusable = true}}, SHIFT(7),
  [39] = {.entry = {.count = 1, .reusable = false}}, SHIFT(7),
  [41] = {.entry = {.count = 1, .reusable = false}}, SHIFT(27),
  [43] = {.entry = {.count = 1, .reusable = false}}, SHIFT(26),
  [45] = {.entry = {.count = 1, .reusable = true}}, SHIFT(18),
  [47] = {.entry = {.count = 1, .reusable = true}}, REDUCE(aux_sym_literal_repeat1, 2),
  [49] = {.entry = {.count = 2, .reusable = true}}, REDUCE(aux_sym_literal_repeat1, 2), SHIFT_REPEAT(7),
  [52] = {.entry = {.count = 2, .reusable = false}}, REDUCE(aux_sym_literal_repeat1, 2), SHIFT_REPEAT(7),
  [55] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym__bracketed_environment_expansion, 4),
  [57] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym__bracketed_environment_expansion, 4),
  [59] = {.entry = {.count = 1, .reusable = true}}, SHIFT(10),
  [61] = {.entry = {.count = 1, .reusable = false}}, SHIFT(6),
  [63] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_variable, 1),
  [65] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym_variable, 1),
  [67] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym__bracketed_environment_expansion, 6),
  [69] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym__bracketed_environment_expansion, 6),
  [71] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym__bracketed_environment_expansion, 5),
  [73] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym__bracketed_environment_expansion, 5),
  [75] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym__mustache_expansion, 5),
  [77] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym__mustache_expansion, 5),
  [79] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym__simple_environment_expansion, 2),
  [81] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym__simple_environment_expansion, 2),
  [83] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym__github_actions_expansion, 2),
  [85] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym__github_actions_expansion, 2),
  [87] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym__mustache_expansion, 3),
  [89] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym__mustache_expansion, 3),
  [91] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym__mustache_expansion, 4),
  [93] = {.entry = {.count = 1, .reusable = false}}, REDUCE(sym__mustache_expansion, 4),
  [95] = {.entry = {.count = 1, .reusable = true}}, SHIFT(28),
  [97] = {.entry = {.count = 1, .reusable = true}}, SHIFT(19),
  [99] = {.entry = {.count = 1, .reusable = true}}, SHIFT(32),
  [101] = {.entry = {.count = 1, .reusable = true}}, SHIFT(12),
  [103] = {.entry = {.count = 1, .reusable = true}}, SHIFT(30),
  [105] = {.entry = {.count = 1, .reusable = true}}, SHIFT(8),
  [107] = {.entry = {.count = 1, .reusable = true}}, SHIFT(29),
  [109] = {.entry = {.count = 1, .reusable = true}}, SHIFT(17),
  [111] = {.entry = {.count = 1, .reusable = true}}, REDUCE(sym_unknown, 1),
  [113] = {.entry = {.count = 1, .reusable = true}}, SHIFT(14),
  [115] = {.entry = {.count = 1, .reusable = true}},  ACCEPT_INPUT(),
  [117] = {.entry = {.count = 1, .reusable = true}}, SHIFT(11),
};

#ifdef __cplusplus
extern "C" {
#endif
#ifdef _WIN32
#define extern __declspec(dllexport)
#endif

extern const TSLanguage *tree_sitter_config_variables(void) {
  static const TSLanguage language = {
    .version = LANGUAGE_VERSION,
    .symbol_count = SYMBOL_COUNT,
    .alias_count = ALIAS_COUNT,
    .token_count = TOKEN_COUNT,
    .external_token_count = EXTERNAL_TOKEN_COUNT,
    .state_count = STATE_COUNT,
    .large_state_count = LARGE_STATE_COUNT,
    .production_id_count = PRODUCTION_ID_COUNT,
    .field_count = FIELD_COUNT,
    .max_alias_sequence_length = MAX_ALIAS_SEQUENCE_LENGTH,
    .parse_table = &ts_parse_table[0][0],
    .small_parse_table = ts_small_parse_table,
    .small_parse_table_map = ts_small_parse_table_map,
    .parse_actions = ts_parse_actions,
    .symbol_names = ts_symbol_names,
    .symbol_metadata = ts_symbol_metadata,
    .public_symbol_map = ts_symbol_map,
    .alias_map = ts_non_terminal_alias_map,
    .alias_sequences = &ts_alias_sequences[0][0],
    .lex_modes = ts_lex_modes,
    .lex_fn = ts_lex,
    .primary_state_ids = ts_primary_state_ids,
  };
  return &language;
}
#ifdef __cplusplus
}
#endif
