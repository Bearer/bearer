#include <tree_sitter/parser.h>

#if defined(__GNUC__) || defined(__clang__)
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wmissing-field-initializers"
#endif

#define LANGUAGE_VERSION 14
#define STATE_COUNT 140
#define LARGE_STATE_COUNT 6
#define SYMBOL_COUNT 67
#define ALIAS_COUNT 0
#define TOKEN_COUNT 42
#define EXTERNAL_TOKEN_COUNT 9
#define FIELD_COUNT 0
#define MAX_ALIAS_SEQUENCE_LENGTH 4
#define PRODUCTION_ID_COUNT 2

enum
{
  anon_sym_LT_BANG = 1,
  aux_sym_doctype_token1 = 2,
  anon_sym_GT = 3,
  sym__doctype = 4,
  anon_sym_LT = 5,
  anon_sym_SLASH_GT = 6,
  anon_sym_LT_SLASH = 7,
  anon_sym_EQ = 8,
  sym_attribute_name = 9,
  sym_attribute_value = 10,
  sym_entity = 11,
  anon_sym_SQUOTE = 12,
  aux_sym_quoted_attribute_value_token1 = 13,
  anon_sym_DQUOTE = 14,
  aux_sym_quoted_attribute_value_token2 = 15,
  sym_text = 16,
  aux_sym_template_code_token1 = 17,
  anon_sym_PERCENT_PERCENT_GT = 18,
  anon_sym_LT_PERCENT = 19,
  anon_sym_LT_PERCENT_ = 20,
  anon_sym_PERCENT_GT = 21,
  anon_sym_DASH_PERCENT_GT = 22,
  anon_sym__PERCENT_GT = 23,
  anon_sym_LT_PERCENT_EQ = 24,
  anon_sym_LT_PERCENT_DASH = 25,
  anon_sym_EQ_PERCENT_GT = 26,
  anon_sym_LT_PERCENT_POUND = 27,
  anon_sym_LT_PERCENTgraphql = 28,
  anon_sym_LBRACE_POUND = 29,
  aux_sym_templating_liquid_comment_token1 = 30,
  anon_sym_POUND_RBRACE = 31,
  anon_sym_LBRACE_PERCENT = 32,
  anon_sym_PERCENT_RBRACE = 33,
  sym__start_tag_name = 34,
  sym__script_start_tag_name = 35,
  sym__style_start_tag_name = 36,
  sym__end_tag_name = 37,
  sym_erroneous_end_tag_name = 38,
  sym__implicit_end_tag = 39,
  sym_raw_text = 40,
  sym_comment = 41,
  sym_fragment = 42,
  sym_doctype = 43,
  sym__node = 44,
  sym_element = 45,
  sym_script_element = 46,
  sym_style_element = 47,
  sym_start_tag = 48,
  sym_script_start_tag = 49,
  sym_style_start_tag = 50,
  sym_self_closing_tag = 51,
  sym_end_tag = 52,
  sym_erroneous_end_tag = 53,
  sym_attribute = 54,
  sym_quoted_attribute_value = 55,
  sym_template = 56,
  sym_template_code = 57,
  sym_template_directive = 58,
  sym_template_output_directive = 59,
  sym_template_comment_directive = 60,
  sym_template_graphql_directive = 61,
  sym_templating_liquid_comment = 62,
  sym_templating_liquid_block = 63,
  aux_sym_fragment_repeat1 = 64,
  aux_sym_start_tag_repeat1 = 65,
  aux_sym_template_code_repeat1 = 66,
};

static const char *const ts_symbol_names[] = {
    [ts_builtin_sym_end] = "end",
    [anon_sym_LT_BANG] = "<!",
    [aux_sym_doctype_token1] = "doctype_token1",
    [anon_sym_GT] = ">",
    [sym__doctype] = "doctype",
    [anon_sym_LT] = "<",
    [anon_sym_SLASH_GT] = "/>",
    [anon_sym_LT_SLASH] = "</",
    [anon_sym_EQ] = "=",
    [sym_attribute_name] = "attribute_name",
    [sym_attribute_value] = "attribute_value",
    [sym_entity] = "entity",
    [anon_sym_SQUOTE] = "'",
    [aux_sym_quoted_attribute_value_token1] = "attribute_value",
    [anon_sym_DQUOTE] = "\"",
    [aux_sym_quoted_attribute_value_token2] = "attribute_value",
    [sym_text] = "text",
    [aux_sym_template_code_token1] = "template_code_token1",
    [anon_sym_PERCENT_PERCENT_GT] = "%%>",
    [anon_sym_LT_PERCENT] = "<%",
    [anon_sym_LT_PERCENT_] = "<%_",
    [anon_sym_PERCENT_GT] = "%>",
    [anon_sym_DASH_PERCENT_GT] = "-%>",
    [anon_sym__PERCENT_GT] = "_%>",
    [anon_sym_LT_PERCENT_EQ] = "<%=",
    [anon_sym_LT_PERCENT_DASH] = "<%-",
    [anon_sym_EQ_PERCENT_GT] = "=%>",
    [anon_sym_LT_PERCENT_POUND] = "<%#",
    [anon_sym_LT_PERCENTgraphql] = "<%graphql",
    [anon_sym_LBRACE_POUND] = "{#",
    [aux_sym_templating_liquid_comment_token1] = "templating_liquid_comment_token1",
    [anon_sym_POUND_RBRACE] = "#}",
    [anon_sym_LBRACE_PERCENT] = "{%",
    [anon_sym_PERCENT_RBRACE] = "%}",
    [sym__start_tag_name] = "tag_name",
    [sym__script_start_tag_name] = "tag_name",
    [sym__style_start_tag_name] = "tag_name",
    [sym__end_tag_name] = "tag_name",
    [sym_erroneous_end_tag_name] = "erroneous_end_tag_name",
    [sym__implicit_end_tag] = "_implicit_end_tag",
    [sym_raw_text] = "raw_text",
    [sym_comment] = "comment",
    [sym_fragment] = "fragment",
    [sym_doctype] = "doctype",
    [sym__node] = "_node",
    [sym_element] = "element",
    [sym_script_element] = "script_element",
    [sym_style_element] = "style_element",
    [sym_start_tag] = "start_tag",
    [sym_script_start_tag] = "start_tag",
    [sym_style_start_tag] = "start_tag",
    [sym_self_closing_tag] = "self_closing_tag",
    [sym_end_tag] = "end_tag",
    [sym_erroneous_end_tag] = "erroneous_end_tag",
    [sym_attribute] = "attribute",
    [sym_quoted_attribute_value] = "quoted_attribute_value",
    [sym_template] = "template",
    [sym_template_code] = "template_code",
    [sym_template_directive] = "template_directive",
    [sym_template_output_directive] = "template_output_directive",
    [sym_template_comment_directive] = "template_comment_directive",
    [sym_template_graphql_directive] = "template_graphql_directive",
    [sym_templating_liquid_comment] = "templating_liquid_comment",
    [sym_templating_liquid_block] = "templating_liquid_block",
    [aux_sym_fragment_repeat1] = "fragment_repeat1",
    [aux_sym_start_tag_repeat1] = "start_tag_repeat1",
    [aux_sym_template_code_repeat1] = "template_code_repeat1",
};

static const TSSymbol ts_symbol_map[] = {
    [ts_builtin_sym_end] = ts_builtin_sym_end,
    [anon_sym_LT_BANG] = anon_sym_LT_BANG,
    [aux_sym_doctype_token1] = aux_sym_doctype_token1,
    [anon_sym_GT] = anon_sym_GT,
    [sym__doctype] = sym__doctype,
    [anon_sym_LT] = anon_sym_LT,
    [anon_sym_SLASH_GT] = anon_sym_SLASH_GT,
    [anon_sym_LT_SLASH] = anon_sym_LT_SLASH,
    [anon_sym_EQ] = anon_sym_EQ,
    [sym_attribute_name] = sym_attribute_name,
    [sym_attribute_value] = sym_attribute_value,
    [sym_entity] = sym_entity,
    [anon_sym_SQUOTE] = anon_sym_SQUOTE,
    [aux_sym_quoted_attribute_value_token1] = sym_attribute_value,
    [anon_sym_DQUOTE] = anon_sym_DQUOTE,
    [aux_sym_quoted_attribute_value_token2] = sym_attribute_value,
    [sym_text] = sym_text,
    [aux_sym_template_code_token1] = aux_sym_template_code_token1,
    [anon_sym_PERCENT_PERCENT_GT] = anon_sym_PERCENT_PERCENT_GT,
    [anon_sym_LT_PERCENT] = anon_sym_LT_PERCENT,
    [anon_sym_LT_PERCENT_] = anon_sym_LT_PERCENT_,
    [anon_sym_PERCENT_GT] = anon_sym_PERCENT_GT,
    [anon_sym_DASH_PERCENT_GT] = anon_sym_DASH_PERCENT_GT,
    [anon_sym__PERCENT_GT] = anon_sym__PERCENT_GT,
    [anon_sym_LT_PERCENT_EQ] = anon_sym_LT_PERCENT_EQ,
    [anon_sym_LT_PERCENT_DASH] = anon_sym_LT_PERCENT_DASH,
    [anon_sym_EQ_PERCENT_GT] = anon_sym_EQ_PERCENT_GT,
    [anon_sym_LT_PERCENT_POUND] = anon_sym_LT_PERCENT_POUND,
    [anon_sym_LT_PERCENTgraphql] = anon_sym_LT_PERCENTgraphql,
    [anon_sym_LBRACE_POUND] = anon_sym_LBRACE_POUND,
    [aux_sym_templating_liquid_comment_token1] = aux_sym_templating_liquid_comment_token1,
    [anon_sym_POUND_RBRACE] = anon_sym_POUND_RBRACE,
    [anon_sym_LBRACE_PERCENT] = anon_sym_LBRACE_PERCENT,
    [anon_sym_PERCENT_RBRACE] = anon_sym_PERCENT_RBRACE,
    [sym__start_tag_name] = sym__start_tag_name,
    [sym__script_start_tag_name] = sym__start_tag_name,
    [sym__style_start_tag_name] = sym__start_tag_name,
    [sym__end_tag_name] = sym__start_tag_name,
    [sym_erroneous_end_tag_name] = sym_erroneous_end_tag_name,
    [sym__implicit_end_tag] = sym__implicit_end_tag,
    [sym_raw_text] = sym_raw_text,
    [sym_comment] = sym_comment,
    [sym_fragment] = sym_fragment,
    [sym_doctype] = sym_doctype,
    [sym__node] = sym__node,
    [sym_element] = sym_element,
    [sym_script_element] = sym_script_element,
    [sym_style_element] = sym_style_element,
    [sym_start_tag] = sym_start_tag,
    [sym_script_start_tag] = sym_start_tag,
    [sym_style_start_tag] = sym_start_tag,
    [sym_self_closing_tag] = sym_self_closing_tag,
    [sym_end_tag] = sym_end_tag,
    [sym_erroneous_end_tag] = sym_erroneous_end_tag,
    [sym_attribute] = sym_attribute,
    [sym_quoted_attribute_value] = sym_quoted_attribute_value,
    [sym_template] = sym_template,
    [sym_template_code] = sym_template_code,
    [sym_template_directive] = sym_template_directive,
    [sym_template_output_directive] = sym_template_output_directive,
    [sym_template_comment_directive] = sym_template_comment_directive,
    [sym_template_graphql_directive] = sym_template_graphql_directive,
    [sym_templating_liquid_comment] = sym_templating_liquid_comment,
    [sym_templating_liquid_block] = sym_templating_liquid_block,
    [aux_sym_fragment_repeat1] = aux_sym_fragment_repeat1,
    [aux_sym_start_tag_repeat1] = aux_sym_start_tag_repeat1,
    [aux_sym_template_code_repeat1] = aux_sym_template_code_repeat1,
};

static const TSSymbolMetadata ts_symbol_metadata[] = {
    [ts_builtin_sym_end] = {
        .visible = false,
        .named = true,
    },
    [anon_sym_LT_BANG] = {
        .visible = true,
        .named = false,
    },
    [aux_sym_doctype_token1] = {
        .visible = false,
        .named = false,
    },
    [anon_sym_GT] = {
        .visible = true,
        .named = false,
    },
    [sym__doctype] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_LT] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_SLASH_GT] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_LT_SLASH] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_EQ] = {
        .visible = true,
        .named = false,
    },
    [sym_attribute_name] = {
        .visible = true,
        .named = true,
    },
    [sym_attribute_value] = {
        .visible = true,
        .named = true,
    },
    [sym_entity] = {
        .visible = true,
        .named = true,
    },
    [anon_sym_SQUOTE] = {
        .visible = true,
        .named = false,
    },
    [aux_sym_quoted_attribute_value_token1] = {
        .visible = true,
        .named = true,
    },
    [anon_sym_DQUOTE] = {
        .visible = true,
        .named = false,
    },
    [aux_sym_quoted_attribute_value_token2] = {
        .visible = true,
        .named = true,
    },
    [sym_text] = {
        .visible = true,
        .named = true,
    },
    [aux_sym_template_code_token1] = {
        .visible = false,
        .named = false,
    },
    [anon_sym_PERCENT_PERCENT_GT] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_LT_PERCENT] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_LT_PERCENT_] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_PERCENT_GT] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_DASH_PERCENT_GT] = {
        .visible = true,
        .named = false,
    },
    [anon_sym__PERCENT_GT] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_LT_PERCENT_EQ] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_LT_PERCENT_DASH] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_EQ_PERCENT_GT] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_LT_PERCENT_POUND] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_LT_PERCENTgraphql] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_LBRACE_POUND] = {
        .visible = true,
        .named = false,
    },
    [aux_sym_templating_liquid_comment_token1] = {
        .visible = false,
        .named = false,
    },
    [anon_sym_POUND_RBRACE] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_LBRACE_PERCENT] = {
        .visible = true,
        .named = false,
    },
    [anon_sym_PERCENT_RBRACE] = {
        .visible = true,
        .named = false,
    },
    [sym__start_tag_name] = {
        .visible = true,
        .named = true,
    },
    [sym__script_start_tag_name] = {
        .visible = true,
        .named = true,
    },
    [sym__style_start_tag_name] = {
        .visible = true,
        .named = true,
    },
    [sym__end_tag_name] = {
        .visible = true,
        .named = true,
    },
    [sym_erroneous_end_tag_name] = {
        .visible = true,
        .named = true,
    },
    [sym__implicit_end_tag] = {
        .visible = false,
        .named = true,
    },
    [sym_raw_text] = {
        .visible = true,
        .named = true,
    },
    [sym_comment] = {
        .visible = true,
        .named = true,
    },
    [sym_fragment] = {
        .visible = true,
        .named = true,
    },
    [sym_doctype] = {
        .visible = true,
        .named = true,
    },
    [sym__node] = {
        .visible = false,
        .named = true,
    },
    [sym_element] = {
        .visible = true,
        .named = true,
    },
    [sym_script_element] = {
        .visible = true,
        .named = true,
    },
    [sym_style_element] = {
        .visible = true,
        .named = true,
    },
    [sym_start_tag] = {
        .visible = true,
        .named = true,
    },
    [sym_script_start_tag] = {
        .visible = true,
        .named = true,
    },
    [sym_style_start_tag] = {
        .visible = true,
        .named = true,
    },
    [sym_self_closing_tag] = {
        .visible = true,
        .named = true,
    },
    [sym_end_tag] = {
        .visible = true,
        .named = true,
    },
    [sym_erroneous_end_tag] = {
        .visible = true,
        .named = true,
    },
    [sym_attribute] = {
        .visible = true,
        .named = true,
    },
    [sym_quoted_attribute_value] = {
        .visible = true,
        .named = true,
    },
    [sym_template] = {
        .visible = true,
        .named = true,
    },
    [sym_template_code] = {
        .visible = true,
        .named = true,
    },
    [sym_template_directive] = {
        .visible = true,
        .named = true,
    },
    [sym_template_output_directive] = {
        .visible = true,
        .named = true,
    },
    [sym_template_comment_directive] = {
        .visible = true,
        .named = true,
    },
    [sym_template_graphql_directive] = {
        .visible = true,
        .named = true,
    },
    [sym_templating_liquid_comment] = {
        .visible = true,
        .named = true,
    },
    [sym_templating_liquid_block] = {
        .visible = true,
        .named = true,
    },
    [aux_sym_fragment_repeat1] = {
        .visible = false,
        .named = false,
    },
    [aux_sym_start_tag_repeat1] = {
        .visible = false,
        .named = false,
    },
    [aux_sym_template_code_repeat1] = {
        .visible = false,
        .named = false,
    },
};

static const TSSymbol ts_alias_sequences[PRODUCTION_ID_COUNT][MAX_ALIAS_SEQUENCE_LENGTH] = {
    [0] = {0},
    [1] = {
        [1] = sym_comment,
    },
};

static const uint16_t ts_non_terminal_alias_map[] = {
    sym_template_code,
    2,
    sym_template_code,
    sym_comment,
    0,
};

static const TSStateId ts_primary_state_ids[STATE_COUNT] = {
    [0] = 0,
    [1] = 1,
    [2] = 2,
    [3] = 2,
    [4] = 4,
    [5] = 4,
    [6] = 6,
    [7] = 7,
    [8] = 6,
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
    [23] = 10,
    [24] = 24,
    [25] = 9,
    [26] = 26,
    [27] = 26,
    [28] = 22,
    [29] = 29,
    [30] = 18,
    [31] = 31,
    [32] = 17,
    [33] = 33,
    [34] = 34,
    [35] = 31,
    [36] = 34,
    [37] = 33,
    [38] = 21,
    [39] = 29,
    [40] = 40,
    [41] = 41,
    [42] = 16,
    [43] = 20,
    [44] = 11,
    [45] = 24,
    [46] = 13,
    [47] = 40,
    [48] = 15,
    [49] = 14,
    [50] = 12,
    [51] = 51,
    [52] = 52,
    [53] = 52,
    [54] = 51,
    [55] = 55,
    [56] = 56,
    [57] = 56,
    [58] = 58,
    [59] = 59,
    [60] = 59,
    [61] = 55,
    [62] = 52,
    [63] = 63,
    [64] = 58,
    [65] = 65,
    [66] = 66,
    [67] = 67,
    [68] = 68,
    [69] = 69,
    [70] = 70,
    [71] = 71,
    [72] = 51,
    [73] = 69,
    [74] = 66,
    [75] = 71,
    [76] = 65,
    [77] = 77,
    [78] = 78,
    [79] = 79,
    [80] = 80,
    [81] = 68,
    [82] = 82,
    [83] = 83,
    [84] = 84,
    [85] = 85,
    [86] = 84,
    [87] = 79,
    [88] = 82,
    [89] = 83,
    [90] = 90,
    [91] = 85,
    [92] = 92,
    [93] = 93,
    [94] = 94,
    [95] = 93,
    [96] = 96,
    [97] = 96,
    [98] = 98,
    [99] = 99,
    [100] = 100,
    [101] = 80,
    [102] = 78,
    [103] = 90,
    [104] = 104,
    [105] = 105,
    [106] = 92,
    [107] = 100,
    [108] = 94,
    [109] = 109,
    [110] = 110,
    [111] = 111,
    [112] = 112,
    [113] = 113,
    [114] = 114,
    [115] = 113,
    [116] = 112,
    [117] = 110,
    [118] = 114,
    [119] = 119,
    [120] = 120,
    [121] = 119,
    [122] = 111,
    [123] = 120,
    [124] = 124,
    [125] = 125,
    [126] = 126,
    [127] = 126,
    [128] = 128,
    [129] = 129,
    [130] = 130,
    [131] = 131,
    [132] = 132,
    [133] = 130,
    [134] = 132,
    [135] = 131,
    [136] = 124,
    [137] = 128,
    [138] = 109,
    [139] = 125,
};

static bool ts_lex(TSLexer *lexer, TSStateId state)
{
  START_LEXER();
  eof = lexer->eof(lexer);
  switch (state)
  {
  case 0:
    if (eof)
      ADVANCE(80);
    if (lookahead == '"')
      ADVANCE(96);
    if (lookahead == '#')
      ADVANCE(67);
    if (lookahead == '%')
      ADVANCE(8);
    if (lookahead == '&')
      ADVANCE(3);
    if (lookahead == '\'')
      ADVANCE(93);
    if (lookahead == '-')
      ADVANCE(9);
    if (lookahead == '/')
      ADVANCE(56);
    if (lookahead == '<')
      ADVANCE(86);
    if (lookahead == '=')
      ADVANCE(89);
    if (lookahead == '>')
      ADVANCE(84);
    if (lookahead == '_')
      ADVANCE(10);
    if (lookahead == '{')
      ADVANCE(4);
    if (lookahead == 'D' ||
        lookahead == 'd')
      ADVANCE(70);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      SKIP(0)
    END_STATE();
  case 1:
    if (lookahead == '"')
      ADVANCE(96);
    if (lookahead == '\'')
      ADVANCE(93);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      SKIP(1)
    if (lookahead != 0 &&
        (lookahead < '<' || '>' < lookahead))
      ADVANCE(91);
    END_STATE();
  case 2:
    if (lookahead == '"')
      ADVANCE(96);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(97);
    if (lookahead != 0)
      ADVANCE(98);
    END_STATE();
  case 3:
    if (lookahead == '#')
      ADVANCE(73);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(54);
    END_STATE();
  case 4:
    if (lookahead == '#')
      ADVANCE(122);
    if (lookahead == '%')
      ADVANCE(126);
    END_STATE();
  case 5:
    if (lookahead == '%')
      ADVANCE(105);
    if (lookahead == '-')
      ADVANCE(106);
    if (lookahead == '=')
      ADVANCE(100);
    if (lookahead == '_')
      ADVANCE(107);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(101);
    if (lookahead != 0)
      ADVANCE(110);
    END_STATE();
  case 6:
    if (lookahead == '%')
      ADVANCE(105);
    if (lookahead == '-')
      ADVANCE(106);
    if (lookahead == '=')
      ADVANCE(108);
    if (lookahead == '_')
      ADVANCE(100);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(102);
    if (lookahead != 0)
      ADVANCE(110);
    END_STATE();
  case 7:
    if (lookahead == '%')
      ADVANCE(105);
    if (lookahead == '-' ||
        lookahead == '=' ||
        lookahead == '_')
      ADVANCE(100);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(103);
    if (lookahead != 0)
      ADVANCE(110);
    END_STATE();
  case 8:
    if (lookahead == '%')
      ADVANCE(57);
    if (lookahead == '>')
      ADVANCE(114);
    if (lookahead == '}')
      ADVANCE(127);
    END_STATE();
  case 9:
    if (lookahead == '%')
      ADVANCE(58);
    END_STATE();
  case 10:
    if (lookahead == '%')
      ADVANCE(59);
    END_STATE();
  case 11:
    if (lookahead == '%')
      ADVANCE(60);
    END_STATE();
  case 12:
    if (lookahead == '%')
      ADVANCE(55);
    if (lookahead == '-')
      ADVANCE(9);
    if (lookahead == '=')
      ADVANCE(11);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      SKIP(12)
    END_STATE();
  case 13:
    if (lookahead == '%')
      ADVANCE(104);
    if (lookahead == '-' ||
        lookahead == '=' ||
        lookahead == '_')
      ADVANCE(100);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(109);
    if (lookahead != 0)
      ADVANCE(110);
    END_STATE();
  case 14:
    if (lookahead == '\'')
      ADVANCE(93);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(94);
    if (lookahead != 0)
      ADVANCE(95);
    END_STATE();
  case 15:
    if (lookahead == '/')
      ADVANCE(56);
    if (lookahead == '=')
      ADVANCE(89);
    if (lookahead == '>')
      ADVANCE(84);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      SKIP(15)
    if (lookahead != 0 &&
        lookahead != '"' &&
        lookahead != '\'' &&
        lookahead != '<')
      ADVANCE(90);
    END_STATE();
  case 16:
    if (lookahead == ';')
      ADVANCE(92);
    END_STATE();
  case 17:
    if (lookahead == ';')
      ADVANCE(92);
    if (('0' <= lookahead && lookahead <= '9'))
      ADVANCE(16);
    END_STATE();
  case 18:
    if (lookahead == ';')
      ADVANCE(92);
    if (('0' <= lookahead && lookahead <= '9'))
      ADVANCE(17);
    END_STATE();
  case 19:
    if (lookahead == ';')
      ADVANCE(92);
    if (('0' <= lookahead && lookahead <= '9'))
      ADVANCE(18);
    END_STATE();
  case 20:
    if (lookahead == ';')
      ADVANCE(92);
    if (('0' <= lookahead && lookahead <= '9'))
      ADVANCE(19);
    END_STATE();
  case 21:
    if (lookahead == ';')
      ADVANCE(92);
    if (('0' <= lookahead && lookahead <= '9') ||
        ('A' <= lookahead && lookahead <= 'F') ||
        ('a' <= lookahead && lookahead <= 'f'))
      ADVANCE(16);
    END_STATE();
  case 22:
    if (lookahead == ';')
      ADVANCE(92);
    if (('0' <= lookahead && lookahead <= '9') ||
        ('A' <= lookahead && lookahead <= 'F') ||
        ('a' <= lookahead && lookahead <= 'f'))
      ADVANCE(21);
    END_STATE();
  case 23:
    if (lookahead == ';')
      ADVANCE(92);
    if (('0' <= lookahead && lookahead <= '9') ||
        ('A' <= lookahead && lookahead <= 'F') ||
        ('a' <= lookahead && lookahead <= 'f'))
      ADVANCE(22);
    END_STATE();
  case 24:
    if (lookahead == ';')
      ADVANCE(92);
    if (('0' <= lookahead && lookahead <= '9') ||
        ('A' <= lookahead && lookahead <= 'F') ||
        ('a' <= lookahead && lookahead <= 'f'))
      ADVANCE(23);
    END_STATE();
  case 25:
    if (lookahead == ';')
      ADVANCE(92);
    if (('0' <= lookahead && lookahead <= '9') ||
        ('A' <= lookahead && lookahead <= 'F') ||
        ('a' <= lookahead && lookahead <= 'f'))
      ADVANCE(24);
    END_STATE();
  case 26:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(16);
    END_STATE();
  case 27:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(26);
    END_STATE();
  case 28:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(27);
    END_STATE();
  case 29:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(28);
    END_STATE();
  case 30:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(29);
    END_STATE();
  case 31:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(30);
    END_STATE();
  case 32:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(31);
    END_STATE();
  case 33:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(32);
    END_STATE();
  case 34:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(33);
    END_STATE();
  case 35:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(34);
    END_STATE();
  case 36:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(35);
    END_STATE();
  case 37:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(36);
    END_STATE();
  case 38:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(37);
    END_STATE();
  case 39:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(38);
    END_STATE();
  case 40:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(39);
    END_STATE();
  case 41:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(40);
    END_STATE();
  case 42:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(41);
    END_STATE();
  case 43:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(42);
    END_STATE();
  case 44:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(43);
    END_STATE();
  case 45:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(44);
    END_STATE();
  case 46:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(45);
    END_STATE();
  case 47:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(46);
    END_STATE();
  case 48:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(47);
    END_STATE();
  case 49:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(48);
    END_STATE();
  case 50:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(49);
    END_STATE();
  case 51:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(50);
    END_STATE();
  case 52:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(51);
    END_STATE();
  case 53:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(52);
    END_STATE();
  case 54:
    if (lookahead == ';')
      ADVANCE(92);
    if (('A' <= lookahead && lookahead <= 'Z') ||
        ('a' <= lookahead && lookahead <= 'z'))
      ADVANCE(53);
    END_STATE();
  case 55:
    if (lookahead == '>')
      ADVANCE(114);
    END_STATE();
  case 56:
    if (lookahead == '>')
      ADVANCE(87);
    END_STATE();
  case 57:
    if (lookahead == '>')
      ADVANCE(111);
    END_STATE();
  case 58:
    if (lookahead == '>')
      ADVANCE(115);
    END_STATE();
  case 59:
    if (lookahead == '>')
      ADVANCE(116);
    END_STATE();
  case 60:
    if (lookahead == '>')
      ADVANCE(119);
    END_STATE();
  case 61:
    if (lookahead == 'a')
      ADVANCE(64);
    END_STATE();
  case 62:
    if (lookahead == 'h')
      ADVANCE(65);
    END_STATE();
  case 63:
    if (lookahead == 'l')
      ADVANCE(121);
    END_STATE();
  case 64:
    if (lookahead == 'p')
      ADVANCE(62);
    END_STATE();
  case 65:
    if (lookahead == 'q')
      ADVANCE(63);
    END_STATE();
  case 66:
    if (lookahead == 'r')
      ADVANCE(61);
    END_STATE();
  case 67:
    if (lookahead == '}')
      ADVANCE(125);
    END_STATE();
  case 68:
    if (lookahead == 'C' ||
        lookahead == 'c')
      ADVANCE(72);
    END_STATE();
  case 69:
    if (lookahead == 'E' ||
        lookahead == 'e')
      ADVANCE(85);
    END_STATE();
  case 70:
    if (lookahead == 'O' ||
        lookahead == 'o')
      ADVANCE(68);
    END_STATE();
  case 71:
    if (lookahead == 'P' ||
        lookahead == 'p')
      ADVANCE(69);
    END_STATE();
  case 72:
    if (lookahead == 'T' ||
        lookahead == 't')
      ADVANCE(74);
    END_STATE();
  case 73:
    if (lookahead == 'X' ||
        lookahead == 'x')
      ADVANCE(78);
    if (('0' <= lookahead && lookahead <= '9'))
      ADVANCE(20);
    END_STATE();
  case 74:
    if (lookahead == 'Y' ||
        lookahead == 'y')
      ADVANCE(71);
    END_STATE();
  case 75:
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(75);
    if (lookahead != 0 &&
        lookahead != '&' &&
        lookahead != '<' &&
        lookahead != '>' &&
        lookahead != '{' &&
        lookahead != '}')
      ADVANCE(99);
    END_STATE();
  case 76:
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(82);
    if (lookahead != 0 &&
        lookahead != '>')
      ADVANCE(83);
    END_STATE();
  case 77:
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(123);
    if (lookahead != 0 &&
        lookahead != '#' &&
        lookahead != '%' &&
        lookahead != '{' &&
        lookahead != '}')
      ADVANCE(124);
    END_STATE();
  case 78:
    if (('0' <= lookahead && lookahead <= '9') ||
        ('A' <= lookahead && lookahead <= 'F') ||
        ('a' <= lookahead && lookahead <= 'f'))
      ADVANCE(25);
    END_STATE();
  case 79:
    if (eof)
      ADVANCE(80);
    if (lookahead == '&')
      ADVANCE(3);
    if (lookahead == '<')
      ADVANCE(86);
    if (lookahead == '{')
      ADVANCE(4);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      SKIP(79)
    if (lookahead != 0 &&
        lookahead != '>' &&
        lookahead != '}')
      ADVANCE(99);
    END_STATE();
  case 80:
    ACCEPT_TOKEN(ts_builtin_sym_end);
    END_STATE();
  case 81:
    ACCEPT_TOKEN(anon_sym_LT_BANG);
    END_STATE();
  case 82:
    ACCEPT_TOKEN(aux_sym_doctype_token1);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(82);
    if (lookahead != 0 &&
        lookahead != '>')
      ADVANCE(83);
    END_STATE();
  case 83:
    ACCEPT_TOKEN(aux_sym_doctype_token1);
    if (lookahead != 0 &&
        lookahead != '>')
      ADVANCE(83);
    END_STATE();
  case 84:
    ACCEPT_TOKEN(anon_sym_GT);
    END_STATE();
  case 85:
    ACCEPT_TOKEN(sym__doctype);
    END_STATE();
  case 86:
    ACCEPT_TOKEN(anon_sym_LT);
    if (lookahead == '!')
      ADVANCE(81);
    if (lookahead == '%')
      ADVANCE(112);
    if (lookahead == '/')
      ADVANCE(88);
    END_STATE();
  case 87:
    ACCEPT_TOKEN(anon_sym_SLASH_GT);
    END_STATE();
  case 88:
    ACCEPT_TOKEN(anon_sym_LT_SLASH);
    END_STATE();
  case 89:
    ACCEPT_TOKEN(anon_sym_EQ);
    END_STATE();
  case 90:
    ACCEPT_TOKEN(sym_attribute_name);
    if (lookahead != 0 &&
        lookahead != '\t' &&
        lookahead != '\n' &&
        lookahead != '\r' &&
        lookahead != ' ' &&
        lookahead != '"' &&
        lookahead != '\'' &&
        lookahead != '/' &&
        (lookahead < '<' || '>' < lookahead))
      ADVANCE(90);
    END_STATE();
  case 91:
    ACCEPT_TOKEN(sym_attribute_value);
    if (lookahead != 0 &&
        lookahead != '\t' &&
        lookahead != '\n' &&
        lookahead != '\r' &&
        lookahead != ' ' &&
        lookahead != '"' &&
        lookahead != '\'' &&
        (lookahead < '<' || '>' < lookahead))
      ADVANCE(91);
    END_STATE();
  case 92:
    ACCEPT_TOKEN(sym_entity);
    END_STATE();
  case 93:
    ACCEPT_TOKEN(anon_sym_SQUOTE);
    END_STATE();
  case 94:
    ACCEPT_TOKEN(aux_sym_quoted_attribute_value_token1);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(94);
    if (lookahead != 0 &&
        lookahead != '\'')
      ADVANCE(95);
    END_STATE();
  case 95:
    ACCEPT_TOKEN(aux_sym_quoted_attribute_value_token1);
    if (lookahead != 0 &&
        lookahead != '\'')
      ADVANCE(95);
    END_STATE();
  case 96:
    ACCEPT_TOKEN(anon_sym_DQUOTE);
    END_STATE();
  case 97:
    ACCEPT_TOKEN(aux_sym_quoted_attribute_value_token2);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(97);
    if (lookahead != 0 &&
        lookahead != '"')
      ADVANCE(98);
    END_STATE();
  case 98:
    ACCEPT_TOKEN(aux_sym_quoted_attribute_value_token2);
    if (lookahead != 0 &&
        lookahead != '"')
      ADVANCE(98);
    END_STATE();
  case 99:
    ACCEPT_TOKEN(sym_text);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(75);
    if (lookahead != 0 &&
        lookahead != '&' &&
        lookahead != '<' &&
        lookahead != '>' &&
        lookahead != '{' &&
        lookahead != '}')
      ADVANCE(99);
    END_STATE();
  case 100:
    ACCEPT_TOKEN(aux_sym_template_code_token1);
    END_STATE();
  case 101:
    ACCEPT_TOKEN(aux_sym_template_code_token1);
    if (lookahead == '%')
      ADVANCE(105);
    if (lookahead == '-')
      ADVANCE(106);
    if (lookahead == '=')
      ADVANCE(100);
    if (lookahead == '_')
      ADVANCE(107);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(101);
    if (lookahead != 0)
      ADVANCE(110);
    END_STATE();
  case 102:
    ACCEPT_TOKEN(aux_sym_template_code_token1);
    if (lookahead == '%')
      ADVANCE(105);
    if (lookahead == '-')
      ADVANCE(106);
    if (lookahead == '=')
      ADVANCE(108);
    if (lookahead == '_')
      ADVANCE(100);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(102);
    if (lookahead != 0)
      ADVANCE(110);
    END_STATE();
  case 103:
    ACCEPT_TOKEN(aux_sym_template_code_token1);
    if (lookahead == '%')
      ADVANCE(105);
    if (lookahead == '-' ||
        lookahead == '=' ||
        lookahead == '_')
      ADVANCE(100);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(103);
    if (lookahead != 0)
      ADVANCE(110);
    END_STATE();
  case 104:
    ACCEPT_TOKEN(aux_sym_template_code_token1);
    if (lookahead == '%')
      ADVANCE(57);
    END_STATE();
  case 105:
    ACCEPT_TOKEN(aux_sym_template_code_token1);
    if (lookahead == '%')
      ADVANCE(57);
    if (lookahead == '>')
      ADVANCE(114);
    END_STATE();
  case 106:
    ACCEPT_TOKEN(aux_sym_template_code_token1);
    if (lookahead == '%')
      ADVANCE(58);
    END_STATE();
  case 107:
    ACCEPT_TOKEN(aux_sym_template_code_token1);
    if (lookahead == '%')
      ADVANCE(59);
    END_STATE();
  case 108:
    ACCEPT_TOKEN(aux_sym_template_code_token1);
    if (lookahead == '%')
      ADVANCE(60);
    END_STATE();
  case 109:
    ACCEPT_TOKEN(aux_sym_template_code_token1);
    if (lookahead == '%')
      ADVANCE(104);
    if (lookahead == '-' ||
        lookahead == '=' ||
        lookahead == '_')
      ADVANCE(100);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(109);
    if (lookahead != 0)
      ADVANCE(110);
    END_STATE();
  case 110:
    ACCEPT_TOKEN(aux_sym_template_code_token1);
    if (lookahead != 0 &&
        lookahead != '%' &&
        lookahead != '-' &&
        lookahead != '=' &&
        lookahead != '_')
      ADVANCE(110);
    END_STATE();
  case 111:
    ACCEPT_TOKEN(anon_sym_PERCENT_PERCENT_GT);
    END_STATE();
  case 112:
    ACCEPT_TOKEN(anon_sym_LT_PERCENT);
    if (lookahead == '#')
      ADVANCE(120);
    if (lookahead == '-')
      ADVANCE(118);
    if (lookahead == '=')
      ADVANCE(117);
    if (lookahead == '_')
      ADVANCE(113);
    if (lookahead == 'g')
      ADVANCE(66);
    END_STATE();
  case 113:
    ACCEPT_TOKEN(anon_sym_LT_PERCENT_);
    END_STATE();
  case 114:
    ACCEPT_TOKEN(anon_sym_PERCENT_GT);
    END_STATE();
  case 115:
    ACCEPT_TOKEN(anon_sym_DASH_PERCENT_GT);
    END_STATE();
  case 116:
    ACCEPT_TOKEN(anon_sym__PERCENT_GT);
    END_STATE();
  case 117:
    ACCEPT_TOKEN(anon_sym_LT_PERCENT_EQ);
    END_STATE();
  case 118:
    ACCEPT_TOKEN(anon_sym_LT_PERCENT_DASH);
    END_STATE();
  case 119:
    ACCEPT_TOKEN(anon_sym_EQ_PERCENT_GT);
    END_STATE();
  case 120:
    ACCEPT_TOKEN(anon_sym_LT_PERCENT_POUND);
    END_STATE();
  case 121:
    ACCEPT_TOKEN(anon_sym_LT_PERCENTgraphql);
    END_STATE();
  case 122:
    ACCEPT_TOKEN(anon_sym_LBRACE_POUND);
    END_STATE();
  case 123:
    ACCEPT_TOKEN(aux_sym_templating_liquid_comment_token1);
    if (lookahead == '\t' ||
        lookahead == '\n' ||
        lookahead == '\r' ||
        lookahead == ' ')
      ADVANCE(123);
    if (lookahead != 0 &&
        lookahead != '#' &&
        lookahead != '%' &&
        lookahead != '{' &&
        lookahead != '}')
      ADVANCE(124);
    END_STATE();
  case 124:
    ACCEPT_TOKEN(aux_sym_templating_liquid_comment_token1);
    if (lookahead != 0 &&
        lookahead != '#' &&
        lookahead != '%' &&
        lookahead != '{' &&
        lookahead != '}')
      ADVANCE(124);
    END_STATE();
  case 125:
    ACCEPT_TOKEN(anon_sym_POUND_RBRACE);
    END_STATE();
  case 126:
    ACCEPT_TOKEN(anon_sym_LBRACE_PERCENT);
    END_STATE();
  case 127:
    ACCEPT_TOKEN(anon_sym_PERCENT_RBRACE);
    END_STATE();
  default:
    return false;
  }
}

static const TSLexMode ts_lex_modes[STATE_COUNT] = {
    [0] = {.lex_state = 0, .external_lex_state = 1},
    [1] = {.lex_state = 79, .external_lex_state = 2},
    [2] = {.lex_state = 79, .external_lex_state = 3},
    [3] = {.lex_state = 79, .external_lex_state = 3},
    [4] = {.lex_state = 79, .external_lex_state = 3},
    [5] = {.lex_state = 79, .external_lex_state = 3},
    [6] = {.lex_state = 79, .external_lex_state = 2},
    [7] = {.lex_state = 79, .external_lex_state = 2},
    [8] = {.lex_state = 79, .external_lex_state = 3},
    [9] = {.lex_state = 79, .external_lex_state = 2},
    [10] = {.lex_state = 79, .external_lex_state = 2},
    [11] = {.lex_state = 79, .external_lex_state = 3},
    [12] = {.lex_state = 79, .external_lex_state = 3},
    [13] = {.lex_state = 79, .external_lex_state = 2},
    [14] = {.lex_state = 79, .external_lex_state = 3},
    [15] = {.lex_state = 79, .external_lex_state = 3},
    [16] = {.lex_state = 79, .external_lex_state = 3},
    [17] = {.lex_state = 79, .external_lex_state = 2},
    [18] = {.lex_state = 79, .external_lex_state = 2},
    [19] = {.lex_state = 79, .external_lex_state = 3},
    [20] = {.lex_state = 79, .external_lex_state = 2},
    [21] = {.lex_state = 79, .external_lex_state = 3},
    [22] = {.lex_state = 79, .external_lex_state = 3},
    [23] = {.lex_state = 79, .external_lex_state = 3},
    [24] = {.lex_state = 79, .external_lex_state = 2},
    [25] = {.lex_state = 79, .external_lex_state = 3},
    [26] = {.lex_state = 79, .external_lex_state = 2},
    [27] = {.lex_state = 79, .external_lex_state = 3},
    [28] = {.lex_state = 79, .external_lex_state = 2},
    [29] = {.lex_state = 79, .external_lex_state = 2},
    [30] = {.lex_state = 79, .external_lex_state = 3},
    [31] = {.lex_state = 79, .external_lex_state = 2},
    [32] = {.lex_state = 79, .external_lex_state = 3},
    [33] = {.lex_state = 79, .external_lex_state = 3},
    [34] = {.lex_state = 79, .external_lex_state = 3},
    [35] = {.lex_state = 79, .external_lex_state = 3},
    [36] = {.lex_state = 79, .external_lex_state = 2},
    [37] = {.lex_state = 79, .external_lex_state = 2},
    [38] = {.lex_state = 79, .external_lex_state = 2},
    [39] = {.lex_state = 79, .external_lex_state = 3},
    [40] = {.lex_state = 79, .external_lex_state = 2},
    [41] = {.lex_state = 79, .external_lex_state = 3},
    [42] = {.lex_state = 79, .external_lex_state = 2},
    [43] = {.lex_state = 79, .external_lex_state = 3},
    [44] = {.lex_state = 79, .external_lex_state = 2},
    [45] = {.lex_state = 79, .external_lex_state = 3},
    [46] = {.lex_state = 79, .external_lex_state = 3},
    [47] = {.lex_state = 79, .external_lex_state = 3},
    [48] = {.lex_state = 79, .external_lex_state = 2},
    [49] = {.lex_state = 79, .external_lex_state = 2},
    [50] = {.lex_state = 79, .external_lex_state = 2},
    [51] = {.lex_state = 5, .external_lex_state = 2},
    [52] = {.lex_state = 6, .external_lex_state = 2},
    [53] = {.lex_state = 5, .external_lex_state = 2},
    [54] = {.lex_state = 6, .external_lex_state = 2},
    [55] = {.lex_state = 15, .external_lex_state = 4},
    [56] = {.lex_state = 7, .external_lex_state = 2},
    [57] = {.lex_state = 7, .external_lex_state = 2},
    [58] = {.lex_state = 15, .external_lex_state = 4},
    [59] = {.lex_state = 15, .external_lex_state = 4},
    [60] = {.lex_state = 15, .external_lex_state = 4},
    [61] = {.lex_state = 15, .external_lex_state = 4},
    [62] = {.lex_state = 7, .external_lex_state = 2},
    [63] = {.lex_state = 15, .external_lex_state = 2},
    [64] = {.lex_state = 15, .external_lex_state = 2},
    [65] = {.lex_state = 1, .external_lex_state = 2},
    [66] = {.lex_state = 13, .external_lex_state = 2},
    [67] = {.lex_state = 15, .external_lex_state = 2},
    [68] = {.lex_state = 15, .external_lex_state = 4},
    [69] = {.lex_state = 13, .external_lex_state = 2},
    [70] = {.lex_state = 15, .external_lex_state = 2},
    [71] = {.lex_state = 13, .external_lex_state = 2},
    [72] = {.lex_state = 7, .external_lex_state = 2},
    [73] = {.lex_state = 13, .external_lex_state = 2},
    [74] = {.lex_state = 13, .external_lex_state = 2},
    [75] = {.lex_state = 13, .external_lex_state = 2},
    [76] = {.lex_state = 1, .external_lex_state = 2},
    [77] = {.lex_state = 15, .external_lex_state = 2},
    [78] = {.lex_state = 15, .external_lex_state = 4},
    [79] = {.lex_state = 0, .external_lex_state = 2},
    [80] = {.lex_state = 15, .external_lex_state = 4},
    [81] = {.lex_state = 15, .external_lex_state = 2},
    [82] = {.lex_state = 0, .external_lex_state = 5},
    [83] = {.lex_state = 0, .external_lex_state = 5},
    [84] = {.lex_state = 12, .external_lex_state = 2},
    [85] = {.lex_state = 0, .external_lex_state = 6},
    [86] = {.lex_state = 12, .external_lex_state = 2},
    [87] = {.lex_state = 0, .external_lex_state = 2},
    [88] = {.lex_state = 0, .external_lex_state = 5},
    [89] = {.lex_state = 0, .external_lex_state = 5},
    [90] = {.lex_state = 15, .external_lex_state = 4},
    [91] = {.lex_state = 0, .external_lex_state = 6},
    [92] = {.lex_state = 0, .external_lex_state = 2},
    [93] = {.lex_state = 2, .external_lex_state = 2},
    [94] = {.lex_state = 0, .external_lex_state = 7},
    [95] = {.lex_state = 2, .external_lex_state = 2},
    [96] = {.lex_state = 14, .external_lex_state = 2},
    [97] = {.lex_state = 14, .external_lex_state = 2},
    [98] = {.lex_state = 0, .external_lex_state = 5},
    [99] = {.lex_state = 0, .external_lex_state = 5},
    [100] = {.lex_state = 0, .external_lex_state = 2},
    [101] = {.lex_state = 15, .external_lex_state = 2},
    [102] = {.lex_state = 15, .external_lex_state = 2},
    [103] = {.lex_state = 15, .external_lex_state = 2},
    [104] = {.lex_state = 0, .external_lex_state = 5},
    [105] = {.lex_state = 0, .external_lex_state = 5},
    [106] = {.lex_state = 0, .external_lex_state = 2},
    [107] = {.lex_state = 0, .external_lex_state = 2},
    [108] = {.lex_state = 0, .external_lex_state = 7},
    [109] = {.lex_state = 0, .external_lex_state = 8},
    [110] = {.lex_state = 0, .external_lex_state = 2},
    [111] = {.lex_state = 0, .external_lex_state = 2},
    [112] = {.lex_state = 0, .external_lex_state = 2},
    [113] = {.lex_state = 0, .external_lex_state = 2},
    [114] = {.lex_state = 0, .external_lex_state = 2},
    [115] = {.lex_state = 0, .external_lex_state = 2},
    [116] = {.lex_state = 0, .external_lex_state = 2},
    [117] = {.lex_state = 0, .external_lex_state = 2},
    [118] = {.lex_state = 0, .external_lex_state = 2},
    [119] = {.lex_state = 0, .external_lex_state = 2},
    [120] = {.lex_state = 0, .external_lex_state = 2},
    [121] = {.lex_state = 0, .external_lex_state = 2},
    [122] = {.lex_state = 0, .external_lex_state = 2},
    [123] = {.lex_state = 0, .external_lex_state = 2},
    [124] = {.lex_state = 76, .external_lex_state = 2},
    [125] = {.lex_state = 0, .external_lex_state = 2},
    [126] = {.lex_state = 0, .external_lex_state = 2},
    [127] = {.lex_state = 0, .external_lex_state = 2},
    [128] = {.lex_state = 0, .external_lex_state = 2},
    [129] = {.lex_state = 0, .external_lex_state = 2},
    [130] = {.lex_state = 0, .external_lex_state = 9},
    [131] = {.lex_state = 77, .external_lex_state = 2},
    [132] = {.lex_state = 77, .external_lex_state = 2},
    [133] = {.lex_state = 0, .external_lex_state = 9},
    [134] = {.lex_state = 77, .external_lex_state = 2},
    [135] = {.lex_state = 77, .external_lex_state = 2},
    [136] = {.lex_state = 76, .external_lex_state = 2},
    [137] = {.lex_state = 0, .external_lex_state = 2},
    [138] = {.lex_state = 0, .external_lex_state = 8},
    [139] = {.lex_state = 0, .external_lex_state = 2},
};

enum
{
  ts_external_token__start_tag_name = 0,
  ts_external_token__script_start_tag_name = 1,
  ts_external_token__style_start_tag_name = 2,
  ts_external_token__end_tag_name = 3,
  ts_external_token_erroneous_end_tag_name = 4,
  ts_external_token_SLASH_GT = 5,
  ts_external_token__implicit_end_tag = 6,
  ts_external_token_raw_text = 7,
  ts_external_token_comment = 8,
};

static const TSSymbol ts_external_scanner_symbol_map[EXTERNAL_TOKEN_COUNT] = {
    [ts_external_token__start_tag_name] = sym__start_tag_name,
    [ts_external_token__script_start_tag_name] = sym__script_start_tag_name,
    [ts_external_token__style_start_tag_name] = sym__style_start_tag_name,
    [ts_external_token__end_tag_name] = sym__end_tag_name,
    [ts_external_token_erroneous_end_tag_name] = sym_erroneous_end_tag_name,
    [ts_external_token_SLASH_GT] = anon_sym_SLASH_GT,
    [ts_external_token__implicit_end_tag] = sym__implicit_end_tag,
    [ts_external_token_raw_text] = sym_raw_text,
    [ts_external_token_comment] = sym_comment,
};

static const bool ts_external_scanner_states[10][EXTERNAL_TOKEN_COUNT] = {
    [1] = {
        [ts_external_token__start_tag_name] = true,
        [ts_external_token__script_start_tag_name] = true,
        [ts_external_token__style_start_tag_name] = true,
        [ts_external_token__end_tag_name] = true,
        [ts_external_token_erroneous_end_tag_name] = true,
        [ts_external_token_SLASH_GT] = true,
        [ts_external_token__implicit_end_tag] = true,
        [ts_external_token_raw_text] = true,
        [ts_external_token_comment] = true,
    },
    [2] = {
        [ts_external_token_comment] = true,
    },
    [3] = {
        [ts_external_token__implicit_end_tag] = true,
        [ts_external_token_comment] = true,
    },
    [4] = {
        [ts_external_token_SLASH_GT] = true,
        [ts_external_token_comment] = true,
    },
    [5] = {
        [ts_external_token_raw_text] = true,
        [ts_external_token_comment] = true,
    },
    [6] = {
        [ts_external_token__start_tag_name] = true,
        [ts_external_token__script_start_tag_name] = true,
        [ts_external_token__style_start_tag_name] = true,
        [ts_external_token_comment] = true,
    },
    [7] = {
        [ts_external_token__end_tag_name] = true,
        [ts_external_token_erroneous_end_tag_name] = true,
        [ts_external_token_comment] = true,
    },
    [8] = {
        [ts_external_token__end_tag_name] = true,
        [ts_external_token_comment] = true,
    },
    [9] = {
        [ts_external_token_erroneous_end_tag_name] = true,
        [ts_external_token_comment] = true,
    },
};

static const uint16_t ts_parse_table[LARGE_STATE_COUNT][SYMBOL_COUNT] = {
    [0] = {
        [ts_builtin_sym_end] = ACTIONS(1),
        [anon_sym_LT_BANG] = ACTIONS(1),
        [anon_sym_GT] = ACTIONS(1),
        [sym__doctype] = ACTIONS(1),
        [anon_sym_LT] = ACTIONS(1),
        [anon_sym_SLASH_GT] = ACTIONS(1),
        [anon_sym_LT_SLASH] = ACTIONS(1),
        [anon_sym_EQ] = ACTIONS(1),
        [sym_entity] = ACTIONS(1),
        [anon_sym_SQUOTE] = ACTIONS(1),
        [anon_sym_DQUOTE] = ACTIONS(1),
        [anon_sym_PERCENT_PERCENT_GT] = ACTIONS(1),
        [anon_sym_LT_PERCENT] = ACTIONS(1),
        [anon_sym_LT_PERCENT_] = ACTIONS(1),
        [anon_sym_PERCENT_GT] = ACTIONS(1),
        [anon_sym_DASH_PERCENT_GT] = ACTIONS(1),
        [anon_sym__PERCENT_GT] = ACTIONS(1),
        [anon_sym_LT_PERCENT_EQ] = ACTIONS(1),
        [anon_sym_LT_PERCENT_DASH] = ACTIONS(1),
        [anon_sym_LT_PERCENT_POUND] = ACTIONS(1),
        [anon_sym_LT_PERCENTgraphql] = ACTIONS(1),
        [anon_sym_LBRACE_POUND] = ACTIONS(1),
        [anon_sym_POUND_RBRACE] = ACTIONS(1),
        [anon_sym_LBRACE_PERCENT] = ACTIONS(1),
        [anon_sym_PERCENT_RBRACE] = ACTIONS(1),
        [sym__start_tag_name] = ACTIONS(1),
        [sym__script_start_tag_name] = ACTIONS(1),
        [sym__style_start_tag_name] = ACTIONS(1),
        [sym__end_tag_name] = ACTIONS(1),
        [sym_erroneous_end_tag_name] = ACTIONS(1),
        [sym__implicit_end_tag] = ACTIONS(1),
        [sym_raw_text] = ACTIONS(1),
        [sym_comment] = ACTIONS(3),
    },
    [1] = {
        [sym_fragment] = STATE(129),
        [sym_doctype] = STATE(7),
        [sym__node] = STATE(7),
        [sym_element] = STATE(7),
        [sym_script_element] = STATE(7),
        [sym_style_element] = STATE(7),
        [sym_start_tag] = STATE(2),
        [sym_script_start_tag] = STATE(83),
        [sym_style_start_tag] = STATE(82),
        [sym_self_closing_tag] = STATE(17),
        [sym_erroneous_end_tag] = STATE(7),
        [sym_template] = STATE(7),
        [sym_template_directive] = STATE(18),
        [sym_template_output_directive] = STATE(18),
        [sym_template_comment_directive] = STATE(18),
        [sym_template_graphql_directive] = STATE(18),
        [sym_templating_liquid_comment] = STATE(18),
        [sym_templating_liquid_block] = STATE(18),
        [aux_sym_fragment_repeat1] = STATE(7),
        [ts_builtin_sym_end] = ACTIONS(5),
        [anon_sym_LT_BANG] = ACTIONS(7),
        [anon_sym_LT] = ACTIONS(9),
        [anon_sym_LT_SLASH] = ACTIONS(11),
        [sym_entity] = ACTIONS(13),
        [sym_text] = ACTIONS(13),
        [anon_sym_LT_PERCENT] = ACTIONS(15),
        [anon_sym_LT_PERCENT_] = ACTIONS(17),
        [anon_sym_LT_PERCENT_EQ] = ACTIONS(19),
        [anon_sym_LT_PERCENT_DASH] = ACTIONS(19),
        [anon_sym_LT_PERCENT_POUND] = ACTIONS(21),
        [anon_sym_LT_PERCENTgraphql] = ACTIONS(23),
        [anon_sym_LBRACE_POUND] = ACTIONS(25),
        [anon_sym_LBRACE_PERCENT] = ACTIONS(27),
        [sym_comment] = ACTIONS(3),
    },
    [2] = {
        [sym_doctype] = STATE(4),
        [sym__node] = STATE(4),
        [sym_element] = STATE(4),
        [sym_script_element] = STATE(4),
        [sym_style_element] = STATE(4),
        [sym_start_tag] = STATE(3),
        [sym_script_start_tag] = STATE(89),
        [sym_style_start_tag] = STATE(88),
        [sym_self_closing_tag] = STATE(32),
        [sym_end_tag] = STATE(9),
        [sym_erroneous_end_tag] = STATE(4),
        [sym_template] = STATE(4),
        [sym_template_directive] = STATE(30),
        [sym_template_output_directive] = STATE(30),
        [sym_template_comment_directive] = STATE(30),
        [sym_template_graphql_directive] = STATE(30),
        [sym_templating_liquid_comment] = STATE(30),
        [sym_templating_liquid_block] = STATE(30),
        [aux_sym_fragment_repeat1] = STATE(4),
        [anon_sym_LT_BANG] = ACTIONS(29),
        [anon_sym_LT] = ACTIONS(31),
        [anon_sym_LT_SLASH] = ACTIONS(33),
        [sym_entity] = ACTIONS(35),
        [sym_text] = ACTIONS(35),
        [anon_sym_LT_PERCENT] = ACTIONS(37),
        [anon_sym_LT_PERCENT_] = ACTIONS(39),
        [anon_sym_LT_PERCENT_EQ] = ACTIONS(41),
        [anon_sym_LT_PERCENT_DASH] = ACTIONS(41),
        [anon_sym_LT_PERCENT_POUND] = ACTIONS(43),
        [anon_sym_LT_PERCENTgraphql] = ACTIONS(45),
        [anon_sym_LBRACE_POUND] = ACTIONS(47),
        [anon_sym_LBRACE_PERCENT] = ACTIONS(49),
        [sym__implicit_end_tag] = ACTIONS(51),
        [sym_comment] = ACTIONS(3),
    },
    [3] = {
        [sym_doctype] = STATE(5),
        [sym__node] = STATE(5),
        [sym_element] = STATE(5),
        [sym_script_element] = STATE(5),
        [sym_style_element] = STATE(5),
        [sym_start_tag] = STATE(3),
        [sym_script_start_tag] = STATE(89),
        [sym_style_start_tag] = STATE(88),
        [sym_self_closing_tag] = STATE(32),
        [sym_end_tag] = STATE(25),
        [sym_erroneous_end_tag] = STATE(5),
        [sym_template] = STATE(5),
        [sym_template_directive] = STATE(30),
        [sym_template_output_directive] = STATE(30),
        [sym_template_comment_directive] = STATE(30),
        [sym_template_graphql_directive] = STATE(30),
        [sym_templating_liquid_comment] = STATE(30),
        [sym_templating_liquid_block] = STATE(30),
        [aux_sym_fragment_repeat1] = STATE(5),
        [anon_sym_LT_BANG] = ACTIONS(29),
        [anon_sym_LT] = ACTIONS(31),
        [anon_sym_LT_SLASH] = ACTIONS(53),
        [sym_entity] = ACTIONS(55),
        [sym_text] = ACTIONS(55),
        [anon_sym_LT_PERCENT] = ACTIONS(37),
        [anon_sym_LT_PERCENT_] = ACTIONS(39),
        [anon_sym_LT_PERCENT_EQ] = ACTIONS(41),
        [anon_sym_LT_PERCENT_DASH] = ACTIONS(41),
        [anon_sym_LT_PERCENT_POUND] = ACTIONS(43),
        [anon_sym_LT_PERCENTgraphql] = ACTIONS(45),
        [anon_sym_LBRACE_POUND] = ACTIONS(47),
        [anon_sym_LBRACE_PERCENT] = ACTIONS(49),
        [sym__implicit_end_tag] = ACTIONS(57),
        [sym_comment] = ACTIONS(3),
    },
    [4] = {
        [sym_doctype] = STATE(8),
        [sym__node] = STATE(8),
        [sym_element] = STATE(8),
        [sym_script_element] = STATE(8),
        [sym_style_element] = STATE(8),
        [sym_start_tag] = STATE(3),
        [sym_script_start_tag] = STATE(89),
        [sym_style_start_tag] = STATE(88),
        [sym_self_closing_tag] = STATE(32),
        [sym_end_tag] = STATE(29),
        [sym_erroneous_end_tag] = STATE(8),
        [sym_template] = STATE(8),
        [sym_template_directive] = STATE(30),
        [sym_template_output_directive] = STATE(30),
        [sym_template_comment_directive] = STATE(30),
        [sym_template_graphql_directive] = STATE(30),
        [sym_templating_liquid_comment] = STATE(30),
        [sym_templating_liquid_block] = STATE(30),
        [aux_sym_fragment_repeat1] = STATE(8),
        [anon_sym_LT_BANG] = ACTIONS(29),
        [anon_sym_LT] = ACTIONS(31),
        [anon_sym_LT_SLASH] = ACTIONS(33),
        [sym_entity] = ACTIONS(59),
        [sym_text] = ACTIONS(59),
        [anon_sym_LT_PERCENT] = ACTIONS(37),
        [anon_sym_LT_PERCENT_] = ACTIONS(39),
        [anon_sym_LT_PERCENT_EQ] = ACTIONS(41),
        [anon_sym_LT_PERCENT_DASH] = ACTIONS(41),
        [anon_sym_LT_PERCENT_POUND] = ACTIONS(43),
        [anon_sym_LT_PERCENTgraphql] = ACTIONS(45),
        [anon_sym_LBRACE_POUND] = ACTIONS(47),
        [anon_sym_LBRACE_PERCENT] = ACTIONS(49),
        [sym__implicit_end_tag] = ACTIONS(61),
        [sym_comment] = ACTIONS(3),
    },
    [5] = {
        [sym_doctype] = STATE(8),
        [sym__node] = STATE(8),
        [sym_element] = STATE(8),
        [sym_script_element] = STATE(8),
        [sym_style_element] = STATE(8),
        [sym_start_tag] = STATE(3),
        [sym_script_start_tag] = STATE(89),
        [sym_style_start_tag] = STATE(88),
        [sym_self_closing_tag] = STATE(32),
        [sym_end_tag] = STATE(39),
        [sym_erroneous_end_tag] = STATE(8),
        [sym_template] = STATE(8),
        [sym_template_directive] = STATE(30),
        [sym_template_output_directive] = STATE(30),
        [sym_template_comment_directive] = STATE(30),
        [sym_template_graphql_directive] = STATE(30),
        [sym_templating_liquid_comment] = STATE(30),
        [sym_templating_liquid_block] = STATE(30),
        [aux_sym_fragment_repeat1] = STATE(8),
        [anon_sym_LT_BANG] = ACTIONS(29),
        [anon_sym_LT] = ACTIONS(31),
        [anon_sym_LT_SLASH] = ACTIONS(53),
        [sym_entity] = ACTIONS(59),
        [sym_text] = ACTIONS(59),
        [anon_sym_LT_PERCENT] = ACTIONS(37),
        [anon_sym_LT_PERCENT_] = ACTIONS(39),
        [anon_sym_LT_PERCENT_EQ] = ACTIONS(41),
        [anon_sym_LT_PERCENT_DASH] = ACTIONS(41),
        [anon_sym_LT_PERCENT_POUND] = ACTIONS(43),
        [anon_sym_LT_PERCENTgraphql] = ACTIONS(45),
        [anon_sym_LBRACE_POUND] = ACTIONS(47),
        [anon_sym_LBRACE_PERCENT] = ACTIONS(49),
        [sym__implicit_end_tag] = ACTIONS(63),
        [sym_comment] = ACTIONS(3),
    },
};

static const uint16_t ts_small_parse_table[] = {
    [0] = 19,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(65),
    1,
    ts_builtin_sym_end,
    ACTIONS(67),
    1,
    anon_sym_LT_BANG,
    ACTIONS(70),
    1,
    anon_sym_LT,
    ACTIONS(73),
    1,
    anon_sym_LT_SLASH,
    ACTIONS(79),
    1,
    anon_sym_LT_PERCENT,
    ACTIONS(82),
    1,
    anon_sym_LT_PERCENT_,
    ACTIONS(88),
    1,
    anon_sym_LT_PERCENT_POUND,
    ACTIONS(91),
    1,
    anon_sym_LT_PERCENTgraphql,
    ACTIONS(94),
    1,
    anon_sym_LBRACE_POUND,
    ACTIONS(97),
    1,
    anon_sym_LBRACE_PERCENT,
    STATE(2),
    1,
    sym_start_tag,
    STATE(17),
    1,
    sym_self_closing_tag,
    STATE(82),
    1,
    sym_style_start_tag,
    STATE(83),
    1,
    sym_script_start_tag,
    ACTIONS(76),
    2,
    sym_entity,
    sym_text,
    ACTIONS(85),
    2,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    STATE(18),
    6,
    sym_template_directive,
    sym_template_output_directive,
    sym_template_comment_directive,
    sym_template_graphql_directive,
    sym_templating_liquid_comment,
    sym_templating_liquid_block,
    STATE(6),
    8,
    sym_doctype,
    sym__node,
    sym_element,
    sym_script_element,
    sym_style_element,
    sym_erroneous_end_tag,
    sym_template,
    aux_sym_fragment_repeat1,
    [72] = 19,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(7),
    1,
    anon_sym_LT_BANG,
    ACTIONS(9),
    1,
    anon_sym_LT,
    ACTIONS(11),
    1,
    anon_sym_LT_SLASH,
    ACTIONS(15),
    1,
    anon_sym_LT_PERCENT,
    ACTIONS(17),
    1,
    anon_sym_LT_PERCENT_,
    ACTIONS(21),
    1,
    anon_sym_LT_PERCENT_POUND,
    ACTIONS(23),
    1,
    anon_sym_LT_PERCENTgraphql,
    ACTIONS(25),
    1,
    anon_sym_LBRACE_POUND,
    ACTIONS(27),
    1,
    anon_sym_LBRACE_PERCENT,
    ACTIONS(100),
    1,
    ts_builtin_sym_end,
    STATE(2),
    1,
    sym_start_tag,
    STATE(17),
    1,
    sym_self_closing_tag,
    STATE(82),
    1,
    sym_style_start_tag,
    STATE(83),
    1,
    sym_script_start_tag,
    ACTIONS(19),
    2,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    ACTIONS(102),
    2,
    sym_entity,
    sym_text,
    STATE(18),
    6,
    sym_template_directive,
    sym_template_output_directive,
    sym_template_comment_directive,
    sym_template_graphql_directive,
    sym_templating_liquid_comment,
    sym_templating_liquid_block,
    STATE(6),
    8,
    sym_doctype,
    sym__node,
    sym_element,
    sym_script_element,
    sym_style_element,
    sym_erroneous_end_tag,
    sym_template,
    aux_sym_fragment_repeat1,
    [144] = 19,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(65),
    1,
    sym__implicit_end_tag,
    ACTIONS(104),
    1,
    anon_sym_LT_BANG,
    ACTIONS(107),
    1,
    anon_sym_LT,
    ACTIONS(110),
    1,
    anon_sym_LT_SLASH,
    ACTIONS(116),
    1,
    anon_sym_LT_PERCENT,
    ACTIONS(119),
    1,
    anon_sym_LT_PERCENT_,
    ACTIONS(125),
    1,
    anon_sym_LT_PERCENT_POUND,
    ACTIONS(128),
    1,
    anon_sym_LT_PERCENTgraphql,
    ACTIONS(131),
    1,
    anon_sym_LBRACE_POUND,
    ACTIONS(134),
    1,
    anon_sym_LBRACE_PERCENT,
    STATE(3),
    1,
    sym_start_tag,
    STATE(32),
    1,
    sym_self_closing_tag,
    STATE(88),
    1,
    sym_style_start_tag,
    STATE(89),
    1,
    sym_script_start_tag,
    ACTIONS(113),
    2,
    sym_entity,
    sym_text,
    ACTIONS(122),
    2,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    STATE(30),
    6,
    sym_template_directive,
    sym_template_output_directive,
    sym_template_comment_directive,
    sym_template_graphql_directive,
    sym_templating_liquid_comment,
    sym_templating_liquid_block,
    STATE(8),
    8,
    sym_doctype,
    sym__node,
    sym_element,
    sym_script_element,
    sym_style_element,
    sym_erroneous_end_tag,
    sym_template,
    aux_sym_fragment_repeat1,
    [216] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(139),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(137),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [238] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(143),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(141),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [260] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(147),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(145),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [282] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(151),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(149),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [304] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(155),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(153),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [326] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(159),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(157),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [348] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(163),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(161),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [370] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(167),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(165),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [392] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(171),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(169),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [414] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(175),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(173),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [436] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(179),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(177),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [458] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(183),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(181),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [480] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(187),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(185),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [502] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(191),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(189),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [524] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(143),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(141),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [546] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(195),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(193),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [568] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(139),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(137),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [590] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(199),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(197),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [612] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(199),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(197),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [634] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(191),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(189),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [656] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(203),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(201),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [678] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(175),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(173),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [700] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(207),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(205),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [722] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(171),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(169),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [744] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(211),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(209),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [766] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(215),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(213),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [788] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(207),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(205),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [810] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(215),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(213),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [832] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(211),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(209),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [854] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(187),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(185),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [876] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(203),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(201),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [898] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(219),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(217),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [920] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(223),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(221),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [942] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(167),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(165),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [964] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(183),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(181),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [986] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(147),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(145),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [1008] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(195),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(193),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [1030] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(155),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(153),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [1052] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(219),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(217),
    12,
    sym__implicit_end_tag,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [1074] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(163),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(161),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [1096] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(159),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(157),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [1118] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(151),
    2,
    anon_sym_LT,
    anon_sym_LT_PERCENT,
    ACTIONS(149),
    12,
    ts_builtin_sym_end,
    anon_sym_LT_BANG,
    anon_sym_LT_SLASH,
    sym_entity,
    sym_text,
    anon_sym_LT_PERCENT_,
    anon_sym_LT_PERCENT_EQ,
    anon_sym_LT_PERCENT_DASH,
    anon_sym_LT_PERCENT_POUND,
    anon_sym_LT_PERCENTgraphql,
    anon_sym_LBRACE_POUND,
    anon_sym_LBRACE_PERCENT,
    [1140] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    STATE(51),
    1,
    aux_sym_template_code_repeat1,
    ACTIONS(225),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    ACTIONS(228),
    3,
    anon_sym_PERCENT_GT,
    anon_sym_DASH_PERCENT_GT,
    anon_sym__PERCENT_GT,
    [1156] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    STATE(54),
    1,
    aux_sym_template_code_repeat1,
    ACTIONS(230),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    ACTIONS(232),
    3,
    anon_sym_PERCENT_GT,
    anon_sym_DASH_PERCENT_GT,
    anon_sym_EQ_PERCENT_GT,
    [1172] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    STATE(51),
    1,
    aux_sym_template_code_repeat1,
    ACTIONS(234),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    ACTIONS(232),
    3,
    anon_sym_PERCENT_GT,
    anon_sym_DASH_PERCENT_GT,
    anon_sym__PERCENT_GT,
    [1188] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    STATE(54),
    1,
    aux_sym_template_code_repeat1,
    ACTIONS(236),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    ACTIONS(228),
    3,
    anon_sym_PERCENT_GT,
    anon_sym_DASH_PERCENT_GT,
    anon_sym_EQ_PERCENT_GT,
    [1204] = 5,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(239),
    1,
    anon_sym_GT,
    ACTIONS(241),
    1,
    anon_sym_SLASH_GT,
    ACTIONS(243),
    1,
    sym_attribute_name,
    STATE(58),
    2,
    sym_attribute,
    aux_sym_start_tag_repeat1,
    [1221] = 5,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(247),
    1,
    anon_sym_PERCENT_GT,
    STATE(62),
    1,
    aux_sym_template_code_repeat1,
    STATE(114),
    1,
    sym_template_code,
    ACTIONS(245),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    [1238] = 5,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(249),
    1,
    anon_sym_PERCENT_GT,
    STATE(62),
    1,
    aux_sym_template_code_repeat1,
    STATE(118),
    1,
    sym_template_code,
    ACTIONS(245),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    [1255] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(253),
    1,
    sym_attribute_name,
    ACTIONS(251),
    2,
    anon_sym_GT,
    anon_sym_SLASH_GT,
    STATE(58),
    2,
    sym_attribute,
    aux_sym_start_tag_repeat1,
    [1270] = 5,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(243),
    1,
    sym_attribute_name,
    ACTIONS(256),
    1,
    anon_sym_GT,
    ACTIONS(258),
    1,
    anon_sym_SLASH_GT,
    STATE(61),
    2,
    sym_attribute,
    aux_sym_start_tag_repeat1,
    [1287] = 5,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(243),
    1,
    sym_attribute_name,
    ACTIONS(256),
    1,
    anon_sym_GT,
    ACTIONS(260),
    1,
    anon_sym_SLASH_GT,
    STATE(55),
    2,
    sym_attribute,
    aux_sym_start_tag_repeat1,
    [1304] = 5,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(239),
    1,
    anon_sym_GT,
    ACTIONS(243),
    1,
    sym_attribute_name,
    ACTIONS(262),
    1,
    anon_sym_SLASH_GT,
    STATE(58),
    2,
    sym_attribute,
    aux_sym_start_tag_repeat1,
    [1321] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(232),
    1,
    anon_sym_PERCENT_GT,
    STATE(72),
    1,
    aux_sym_template_code_repeat1,
    ACTIONS(264),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    [1335] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(266),
    1,
    anon_sym_GT,
    ACTIONS(268),
    1,
    sym_attribute_name,
    STATE(70),
    2,
    sym_attribute,
    aux_sym_start_tag_repeat1,
    [1349] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(251),
    1,
    anon_sym_GT,
    ACTIONS(270),
    1,
    sym_attribute_name,
    STATE(64),
    2,
    sym_attribute,
    aux_sym_start_tag_repeat1,
    [1363] = 5,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(273),
    1,
    sym_attribute_value,
    ACTIONS(275),
    1,
    anon_sym_SQUOTE,
    ACTIONS(277),
    1,
    anon_sym_DQUOTE,
    STATE(80),
    1,
    sym_quoted_attribute_value,
    [1379] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    STATE(62),
    1,
    aux_sym_template_code_repeat1,
    STATE(115),
    1,
    sym_template_code,
    ACTIONS(245),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    [1393] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(268),
    1,
    sym_attribute_name,
    ACTIONS(279),
    1,
    anon_sym_GT,
    STATE(64),
    2,
    sym_attribute,
    aux_sym_start_tag_repeat1,
    [1407] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(283),
    1,
    anon_sym_EQ,
    ACTIONS(281),
    3,
    anon_sym_GT,
    anon_sym_SLASH_GT,
    sym_attribute_name,
    [1419] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    STATE(52),
    1,
    aux_sym_template_code_repeat1,
    STATE(86),
    1,
    sym_template_code,
    ACTIONS(285),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    [1433] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(268),
    1,
    sym_attribute_name,
    ACTIONS(287),
    1,
    anon_sym_GT,
    STATE(64),
    2,
    sym_attribute,
    aux_sym_start_tag_repeat1,
    [1447] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    STATE(53),
    1,
    aux_sym_template_code_repeat1,
    STATE(79),
    1,
    sym_template_code,
    ACTIONS(289),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    [1461] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(228),
    1,
    anon_sym_PERCENT_GT,
    STATE(72),
    1,
    aux_sym_template_code_repeat1,
    ACTIONS(291),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    [1475] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    STATE(52),
    1,
    aux_sym_template_code_repeat1,
    STATE(84),
    1,
    sym_template_code,
    ACTIONS(285),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    [1489] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    STATE(62),
    1,
    aux_sym_template_code_repeat1,
    STATE(113),
    1,
    sym_template_code,
    ACTIONS(245),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    [1503] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    STATE(53),
    1,
    aux_sym_template_code_repeat1,
    STATE(87),
    1,
    sym_template_code,
    ACTIONS(289),
    2,
    aux_sym_template_code_token1,
    anon_sym_PERCENT_PERCENT_GT,
    [1517] = 5,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(294),
    1,
    sym_attribute_value,
    ACTIONS(296),
    1,
    anon_sym_SQUOTE,
    ACTIONS(298),
    1,
    anon_sym_DQUOTE,
    STATE(101),
    1,
    sym_quoted_attribute_value,
    [1533] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(268),
    1,
    sym_attribute_name,
    ACTIONS(300),
    1,
    anon_sym_GT,
    STATE(67),
    2,
    sym_attribute,
    aux_sym_start_tag_repeat1,
    [1547] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(302),
    3,
    anon_sym_GT,
    anon_sym_SLASH_GT,
    sym_attribute_name,
    [1556] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(304),
    3,
    anon_sym_PERCENT_GT,
    anon_sym_DASH_PERCENT_GT,
    anon_sym__PERCENT_GT,
    [1565] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(306),
    3,
    anon_sym_GT,
    anon_sym_SLASH_GT,
    sym_attribute_name,
    [1574] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(308),
    1,
    anon_sym_EQ,
    ACTIONS(281),
    2,
    anon_sym_GT,
    sym_attribute_name,
    [1585] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(310),
    1,
    anon_sym_LT_SLASH,
    ACTIONS(312),
    1,
    sym_raw_text,
    STATE(38),
    1,
    sym_end_tag,
    [1598] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(310),
    1,
    anon_sym_LT_SLASH,
    ACTIONS(314),
    1,
    sym_raw_text,
    STATE(10),
    1,
    sym_end_tag,
    [1611] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(316),
    3,
    anon_sym_PERCENT_GT,
    anon_sym_DASH_PERCENT_GT,
    anon_sym_EQ_PERCENT_GT,
    [1620] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(318),
    1,
    sym__start_tag_name,
    ACTIONS(320),
    1,
    sym__script_start_tag_name,
    ACTIONS(322),
    1,
    sym__style_start_tag_name,
    [1633] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(324),
    3,
    anon_sym_PERCENT_GT,
    anon_sym_DASH_PERCENT_GT,
    anon_sym_EQ_PERCENT_GT,
    [1642] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(326),
    3,
    anon_sym_PERCENT_GT,
    anon_sym_DASH_PERCENT_GT,
    anon_sym__PERCENT_GT,
    [1651] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(328),
    1,
    anon_sym_LT_SLASH,
    ACTIONS(330),
    1,
    sym_raw_text,
    STATE(21),
    1,
    sym_end_tag,
    [1664] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(328),
    1,
    anon_sym_LT_SLASH,
    ACTIONS(332),
    1,
    sym_raw_text,
    STATE(23),
    1,
    sym_end_tag,
    [1677] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(334),
    3,
    anon_sym_GT,
    anon_sym_SLASH_GT,
    sym_attribute_name,
    [1686] = 4,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(320),
    1,
    sym__script_start_tag_name,
    ACTIONS(322),
    1,
    sym__style_start_tag_name,
    ACTIONS(336),
    1,
    sym__start_tag_name,
    [1699] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(310),
    1,
    anon_sym_LT_SLASH,
    STATE(28),
    1,
    sym_end_tag,
    [1709] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(338),
    1,
    anon_sym_DQUOTE,
    ACTIONS(340),
    1,
    aux_sym_quoted_attribute_value_token2,
    [1719] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(342),
    1,
    sym__end_tag_name,
    ACTIONS(344),
    1,
    sym_erroneous_end_tag_name,
    [1729] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(346),
    1,
    anon_sym_DQUOTE,
    ACTIONS(348),
    1,
    aux_sym_quoted_attribute_value_token2,
    [1739] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(346),
    1,
    anon_sym_SQUOTE,
    ACTIONS(350),
    1,
    aux_sym_quoted_attribute_value_token1,
    [1749] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(338),
    1,
    anon_sym_SQUOTE,
    ACTIONS(352),
    1,
    aux_sym_quoted_attribute_value_token1,
    [1759] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(354),
    2,
    sym_raw_text,
    anon_sym_LT_SLASH,
    [1767] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(356),
    2,
    sym_raw_text,
    anon_sym_LT_SLASH,
    [1775] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(328),
    1,
    anon_sym_LT_SLASH,
    STATE(45),
    1,
    sym_end_tag,
    [1785] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(306),
    2,
    anon_sym_GT,
    sym_attribute_name,
    [1793] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(302),
    2,
    anon_sym_GT,
    sym_attribute_name,
    [1801] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(334),
    2,
    anon_sym_GT,
    sym_attribute_name,
    [1809] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(358),
    2,
    sym_raw_text,
    anon_sym_LT_SLASH,
    [1817] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(360),
    2,
    sym_raw_text,
    anon_sym_LT_SLASH,
    [1825] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(328),
    1,
    anon_sym_LT_SLASH,
    STATE(22),
    1,
    sym_end_tag,
    [1835] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(310),
    1,
    anon_sym_LT_SLASH,
    STATE(24),
    1,
    sym_end_tag,
    [1845] = 3,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(344),
    1,
    sym_erroneous_end_tag_name,
    ACTIONS(362),
    1,
    sym__end_tag_name,
    [1855] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(362),
    1,
    sym__end_tag_name,
    [1862] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(364),
    1,
    anon_sym_PERCENT_RBRACE,
    [1869] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(366),
    1,
    anon_sym_GT,
    [1876] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(368),
    1,
    anon_sym_POUND_RBRACE,
    [1883] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(370),
    1,
    anon_sym_PERCENT_GT,
    [1890] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(372),
    1,
    anon_sym_PERCENT_GT,
    [1897] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(374),
    1,
    anon_sym_PERCENT_GT,
    [1904] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(376),
    1,
    anon_sym_POUND_RBRACE,
    [1911] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(378),
    1,
    anon_sym_PERCENT_RBRACE,
    [1918] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(380),
    1,
    anon_sym_PERCENT_GT,
    [1925] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(382),
    1,
    anon_sym_GT,
    [1932] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(384),
    1,
    anon_sym_GT,
    [1939] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(386),
    1,
    anon_sym_GT,
    [1946] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(388),
    1,
    anon_sym_GT,
    [1953] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(390),
    1,
    anon_sym_GT,
    [1960] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(392),
    1,
    aux_sym_doctype_token1,
    [1967] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(394),
    1,
    sym__doctype,
    [1974] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(396),
    1,
    anon_sym_SQUOTE,
    [1981] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(398),
    1,
    anon_sym_SQUOTE,
    [1988] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(398),
    1,
    anon_sym_DQUOTE,
    [1995] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(400),
    1,
    ts_builtin_sym_end,
    [2002] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(344),
    1,
    sym_erroneous_end_tag_name,
    [2009] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(402),
    1,
    aux_sym_templating_liquid_comment_token1,
    [2016] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(404),
    1,
    aux_sym_templating_liquid_comment_token1,
    [2023] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(406),
    1,
    sym_erroneous_end_tag_name,
    [2030] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(408),
    1,
    aux_sym_templating_liquid_comment_token1,
    [2037] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(410),
    1,
    aux_sym_templating_liquid_comment_token1,
    [2044] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(412),
    1,
    aux_sym_doctype_token1,
    [2051] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(396),
    1,
    anon_sym_DQUOTE,
    [2058] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(342),
    1,
    sym__end_tag_name,
    [2065] = 2,
    ACTIONS(3),
    1,
    sym_comment,
    ACTIONS(414),
    1,
    sym__doctype,
};

static const uint32_t ts_small_parse_table_map[] = {
    [SMALL_STATE(6)] = 0,
    [SMALL_STATE(7)] = 72,
    [SMALL_STATE(8)] = 144,
    [SMALL_STATE(9)] = 216,
    [SMALL_STATE(10)] = 238,
    [SMALL_STATE(11)] = 260,
    [SMALL_STATE(12)] = 282,
    [SMALL_STATE(13)] = 304,
    [SMALL_STATE(14)] = 326,
    [SMALL_STATE(15)] = 348,
    [SMALL_STATE(16)] = 370,
    [SMALL_STATE(17)] = 392,
    [SMALL_STATE(18)] = 414,
    [SMALL_STATE(19)] = 436,
    [SMALL_STATE(20)] = 458,
    [SMALL_STATE(21)] = 480,
    [SMALL_STATE(22)] = 502,
    [SMALL_STATE(23)] = 524,
    [SMALL_STATE(24)] = 546,
    [SMALL_STATE(25)] = 568,
    [SMALL_STATE(26)] = 590,
    [SMALL_STATE(27)] = 612,
    [SMALL_STATE(28)] = 634,
    [SMALL_STATE(29)] = 656,
    [SMALL_STATE(30)] = 678,
    [SMALL_STATE(31)] = 700,
    [SMALL_STATE(32)] = 722,
    [SMALL_STATE(33)] = 744,
    [SMALL_STATE(34)] = 766,
    [SMALL_STATE(35)] = 788,
    [SMALL_STATE(36)] = 810,
    [SMALL_STATE(37)] = 832,
    [SMALL_STATE(38)] = 854,
    [SMALL_STATE(39)] = 876,
    [SMALL_STATE(40)] = 898,
    [SMALL_STATE(41)] = 920,
    [SMALL_STATE(42)] = 942,
    [SMALL_STATE(43)] = 964,
    [SMALL_STATE(44)] = 986,
    [SMALL_STATE(45)] = 1008,
    [SMALL_STATE(46)] = 1030,
    [SMALL_STATE(47)] = 1052,
    [SMALL_STATE(48)] = 1074,
    [SMALL_STATE(49)] = 1096,
    [SMALL_STATE(50)] = 1118,
    [SMALL_STATE(51)] = 1140,
    [SMALL_STATE(52)] = 1156,
    [SMALL_STATE(53)] = 1172,
    [SMALL_STATE(54)] = 1188,
    [SMALL_STATE(55)] = 1204,
    [SMALL_STATE(56)] = 1221,
    [SMALL_STATE(57)] = 1238,
    [SMALL_STATE(58)] = 1255,
    [SMALL_STATE(59)] = 1270,
    [SMALL_STATE(60)] = 1287,
    [SMALL_STATE(61)] = 1304,
    [SMALL_STATE(62)] = 1321,
    [SMALL_STATE(63)] = 1335,
    [SMALL_STATE(64)] = 1349,
    [SMALL_STATE(65)] = 1363,
    [SMALL_STATE(66)] = 1379,
    [SMALL_STATE(67)] = 1393,
    [SMALL_STATE(68)] = 1407,
    [SMALL_STATE(69)] = 1419,
    [SMALL_STATE(70)] = 1433,
    [SMALL_STATE(71)] = 1447,
    [SMALL_STATE(72)] = 1461,
    [SMALL_STATE(73)] = 1475,
    [SMALL_STATE(74)] = 1489,
    [SMALL_STATE(75)] = 1503,
    [SMALL_STATE(76)] = 1517,
    [SMALL_STATE(77)] = 1533,
    [SMALL_STATE(78)] = 1547,
    [SMALL_STATE(79)] = 1556,
    [SMALL_STATE(80)] = 1565,
    [SMALL_STATE(81)] = 1574,
    [SMALL_STATE(82)] = 1585,
    [SMALL_STATE(83)] = 1598,
    [SMALL_STATE(84)] = 1611,
    [SMALL_STATE(85)] = 1620,
    [SMALL_STATE(86)] = 1633,
    [SMALL_STATE(87)] = 1642,
    [SMALL_STATE(88)] = 1651,
    [SMALL_STATE(89)] = 1664,
    [SMALL_STATE(90)] = 1677,
    [SMALL_STATE(91)] = 1686,
    [SMALL_STATE(92)] = 1699,
    [SMALL_STATE(93)] = 1709,
    [SMALL_STATE(94)] = 1719,
    [SMALL_STATE(95)] = 1729,
    [SMALL_STATE(96)] = 1739,
    [SMALL_STATE(97)] = 1749,
    [SMALL_STATE(98)] = 1759,
    [SMALL_STATE(99)] = 1767,
    [SMALL_STATE(100)] = 1775,
    [SMALL_STATE(101)] = 1785,
    [SMALL_STATE(102)] = 1793,
    [SMALL_STATE(103)] = 1801,
    [SMALL_STATE(104)] = 1809,
    [SMALL_STATE(105)] = 1817,
    [SMALL_STATE(106)] = 1825,
    [SMALL_STATE(107)] = 1835,
    [SMALL_STATE(108)] = 1845,
    [SMALL_STATE(109)] = 1855,
    [SMALL_STATE(110)] = 1862,
    [SMALL_STATE(111)] = 1869,
    [SMALL_STATE(112)] = 1876,
    [SMALL_STATE(113)] = 1883,
    [SMALL_STATE(114)] = 1890,
    [SMALL_STATE(115)] = 1897,
    [SMALL_STATE(116)] = 1904,
    [SMALL_STATE(117)] = 1911,
    [SMALL_STATE(118)] = 1918,
    [SMALL_STATE(119)] = 1925,
    [SMALL_STATE(120)] = 1932,
    [SMALL_STATE(121)] = 1939,
    [SMALL_STATE(122)] = 1946,
    [SMALL_STATE(123)] = 1953,
    [SMALL_STATE(124)] = 1960,
    [SMALL_STATE(125)] = 1967,
    [SMALL_STATE(126)] = 1974,
    [SMALL_STATE(127)] = 1981,
    [SMALL_STATE(128)] = 1988,
    [SMALL_STATE(129)] = 1995,
    [SMALL_STATE(130)] = 2002,
    [SMALL_STATE(131)] = 2009,
    [SMALL_STATE(132)] = 2016,
    [SMALL_STATE(133)] = 2023,
    [SMALL_STATE(134)] = 2030,
    [SMALL_STATE(135)] = 2037,
    [SMALL_STATE(136)] = 2044,
    [SMALL_STATE(137)] = 2051,
    [SMALL_STATE(138)] = 2058,
    [SMALL_STATE(139)] = 2065,
};

static const TSParseActionEntry ts_parse_actions[] = {
    [0] = {.entry = {.count = 0, .reusable = false}},
    [1] = {.entry = {.count = 1, .reusable = false}},
    RECOVER(),
    [3] = {.entry = {.count = 1, .reusable = true}},
    SHIFT_EXTRA(),
    [5] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_fragment, 0),
    [7] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(125),
    [9] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(91),
    [11] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(133),
    [13] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(7),
    [15] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(71),
    [17] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(71),
    [19] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(73),
    [21] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(57),
    [23] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(74),
    [25] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(132),
    [27] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(131),
    [29] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(139),
    [31] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(85),
    [33] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(108),
    [35] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(4),
    [37] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(75),
    [39] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(75),
    [41] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(69),
    [43] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(56),
    [45] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(66),
    [47] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(134),
    [49] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(135),
    [51] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(9),
    [53] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(94),
    [55] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(5),
    [57] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(25),
    [59] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(8),
    [61] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(29),
    [63] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(39),
    [65] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    [67] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(125),
    [70] = {.entry = {.count = 2, .reusable = false}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(91),
    [73] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(133),
    [76] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(6),
    [79] = {.entry = {.count = 2, .reusable = false}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(71),
    [82] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(71),
    [85] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(73),
    [88] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(57),
    [91] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(74),
    [94] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(132),
    [97] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(131),
    [100] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_fragment, 1),
    [102] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(6),
    [104] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(139),
    [107] = {.entry = {.count = 2, .reusable = false}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(85),
    [110] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(130),
    [113] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(8),
    [116] = {.entry = {.count = 2, .reusable = false}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(75),
    [119] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(75),
    [122] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(69),
    [125] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(56),
    [128] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(66),
    [131] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(134),
    [134] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_fragment_repeat1, 2),
    SHIFT_REPEAT(135),
    [137] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_element, 2),
    [139] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_element, 2),
    [141] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_script_element, 2),
    [143] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_script_element, 2),
    [145] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_template_comment_directive, 3, .production_id = 1),
    [147] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_template_comment_directive, 3, .production_id = 1),
    [149] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_template_output_directive, 3),
    [151] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_template_output_directive, 3),
    [153] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_self_closing_tag, 4),
    [155] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_self_closing_tag, 4),
    [157] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_template_directive, 3),
    [159] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_template_directive, 3),
    [161] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_erroneous_end_tag, 3),
    [163] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_erroneous_end_tag, 3),
    [165] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_self_closing_tag, 3),
    [167] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_self_closing_tag, 3),
    [169] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_element, 1),
    [171] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_element, 1),
    [173] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_template, 1),
    [175] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_template, 1),
    [177] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_start_tag, 4),
    [179] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_start_tag, 4),
    [181] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_doctype, 4),
    [183] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_doctype, 4),
    [185] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_style_element, 2),
    [187] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_style_element, 2),
    [189] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_script_element, 3),
    [191] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_script_element, 3),
    [193] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_style_element, 3),
    [195] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_style_element, 3),
    [197] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_template_comment_directive, 2),
    [199] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_template_comment_directive, 2),
    [201] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_element, 3),
    [203] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_element, 3),
    [205] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_templating_liquid_block, 3),
    [207] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_templating_liquid_block, 3),
    [209] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_template_graphql_directive, 3),
    [211] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_template_graphql_directive, 3),
    [213] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_templating_liquid_comment, 3),
    [215] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_templating_liquid_comment, 3),
    [217] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_end_tag, 3),
    [219] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_end_tag, 3),
    [221] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_start_tag, 3),
    [223] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_start_tag, 3),
    [225] = {.entry = {.count = 2, .reusable = false}},
    REDUCE(aux_sym_template_code_repeat1, 2),
    SHIFT_REPEAT(51),
    [228] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(aux_sym_template_code_repeat1, 2),
    [230] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(54),
    [232] = {.entry = {.count = 1, .reusable = false}},
    REDUCE(sym_template_code, 1),
    [234] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(51),
    [236] = {.entry = {.count = 2, .reusable = false}},
    REDUCE(aux_sym_template_code_repeat1, 2),
    SHIFT_REPEAT(54),
    [239] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(19),
    [241] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(13),
    [243] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(68),
    [245] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(62),
    [247] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(27),
    [249] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(26),
    [251] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(aux_sym_start_tag_repeat1, 2),
    [253] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_start_tag_repeat1, 2),
    SHIFT_REPEAT(68),
    [256] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(41),
    [258] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(16),
    [260] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(42),
    [262] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(46),
    [264] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(72),
    [266] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(104),
    [268] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(81),
    [270] = {.entry = {.count = 2, .reusable = true}},
    REDUCE(aux_sym_start_tag_repeat1, 2),
    SHIFT_REPEAT(81),
    [273] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(80),
    [275] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(96),
    [277] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(95),
    [279] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(99),
    [281] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_attribute, 1),
    [283] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(65),
    [285] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(52),
    [287] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(98),
    [289] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(53),
    [291] = {.entry = {.count = 2, .reusable = false}},
    REDUCE(aux_sym_template_code_repeat1, 2),
    SHIFT_REPEAT(72),
    [294] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(101),
    [296] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(97),
    [298] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(93),
    [300] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(105),
    [302] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_quoted_attribute_value, 2),
    [304] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(49),
    [306] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_attribute, 3),
    [308] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(76),
    [310] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(109),
    [312] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(107),
    [314] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(92),
    [316] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(50),
    [318] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(59),
    [320] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(77),
    [322] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(63),
    [324] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(12),
    [326] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(14),
    [328] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(138),
    [330] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(100),
    [332] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(106),
    [334] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_quoted_attribute_value, 3),
    [336] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(60),
    [338] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(102),
    [340] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(128),
    [342] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(123),
    [344] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(111),
    [346] = {.entry = {.count = 1, .reusable = false}},
    SHIFT(78),
    [348] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(137),
    [350] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(126),
    [352] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(127),
    [354] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_style_start_tag, 4),
    [356] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_script_start_tag, 4),
    [358] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_style_start_tag, 3),
    [360] = {.entry = {.count = 1, .reusable = true}},
    REDUCE(sym_script_start_tag, 3),
    [362] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(120),
    [364] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(31),
    [366] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(15),
    [368] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(36),
    [370] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(37),
    [372] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(11),
    [374] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(33),
    [376] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(34),
    [378] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(35),
    [380] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(44),
    [382] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(20),
    [384] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(40),
    [386] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(43),
    [388] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(48),
    [390] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(47),
    [392] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(119),
    [394] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(124),
    [396] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(90),
    [398] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(103),
    [400] = {.entry = {.count = 1, .reusable = true}},
    ACCEPT_INPUT(),
    [402] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(110),
    [404] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(112),
    [406] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(122),
    [408] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(116),
    [410] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(117),
    [412] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(121),
    [414] = {.entry = {.count = 1, .reusable = true}},
    SHIFT(136),
};

#ifdef __cplusplus
extern "C"
{
#endif
  void *tree_sitter_html2_external_scanner_create(void);
  void tree_sitter_html2_external_scanner_destroy(void *);
  bool tree_sitter_html2_external_scanner_scan(void *, TSLexer *, const bool *);
  unsigned tree_sitter_html2_external_scanner_serialize(void *, char *);
  void tree_sitter_html2_external_scanner_deserialize(void *, const char *, unsigned);

#ifdef _WIN32
#define extern __declspec(dllexport)
#endif

  extern const TSLanguage *tree_sitter_html2(void)
  {
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
        .external_scanner = {
            &ts_external_scanner_states[0][0],
            ts_external_scanner_symbol_map,
            tree_sitter_html2_external_scanner_create,
            tree_sitter_html2_external_scanner_destroy,
            tree_sitter_html2_external_scanner_scan,
            tree_sitter_html2_external_scanner_serialize,
            tree_sitter_html2_external_scanner_deserialize,
        },
        .primary_state_ids = ts_primary_state_ids,
    };
    return &language;
  }
#ifdef __cplusplus
}
#endif
