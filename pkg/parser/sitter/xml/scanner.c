#include "scanner.h"
#include "tree_sitter/array.h"

#include <string.h>

typedef Array(char) String;

typedef Array(String) Vector;

static inline bool string_eq(String *a, String *b) {
    if (a->size != b->size) {
        return false;
    }
    return memcmp(a->contents, b->contents, a->size) == 0;
}

static String scan_tag_name(TSLexer *lexer) {
    String tag_name = array_new();
    if (is_valid_name_start_char(lexer->lookahead)) {
        array_push(&tag_name, (char)lexer->lookahead);
        advance(lexer);
    }
    while (is_valid_name_char(lexer->lookahead)) {
        array_push(&tag_name, (char)lexer->lookahead);
        advance(lexer);
    }
    return tag_name;
}

static bool scan_start_tag_name(Vector *tags, TSLexer *lexer) {
    String tag_name = scan_tag_name(lexer);
    if (tag_name.size == 0) {
        array_delete(&tag_name);
        return false;
    }

    lexer->result_symbol = START_TAG_NAME;
    array_push(tags, tag_name);
    return true;
}

static bool scan_end_tag_name(Vector *tags, TSLexer *lexer) {
    String tag_name = scan_tag_name(lexer);
    if (tag_name.size == 0) {
        array_delete(&tag_name);
        return false;
    }

    if (tags->size > 0 && string_eq(array_back(tags), &tag_name)) {
        array_delete(&array_pop(tags));
        lexer->result_symbol = END_TAG_NAME;
    } else {
        lexer->result_symbol = ERRONEOUS_END_NAME;
    }
    array_delete(&tag_name);
    return lexer->result_symbol == END_TAG_NAME;
}

static bool scan_self_closing_tag_delimiter(Vector *tags, TSLexer *lexer) {
    advance(lexer);
    advance_if_eq(lexer, '>');
    if (tags->size > 0) {
        array_delete(&array_pop(tags));
        lexer->result_symbol = SELF_CLOSING_TAG_DELIMITER;
    }
    return true;
}

/// Check if the lexer is in error recovery mode
static inline bool in_error_recovery(const bool *valid_symbols) {
    return valid_symbols[PI_TARGET] && valid_symbols[PI_CONTENT] && valid_symbols[COMMENT] &&
           valid_symbols[CHAR_DATA] && valid_symbols[CDATA];
}

/// Scan for a CharData node
static bool scan_char_data(TSLexer *lexer) {
    bool advanced_once = false;

    while (!lexer->eof(lexer) && lexer->lookahead != '<' && lexer->lookahead != '&') {
        if (lexer->lookahead == ']') {
            lexer->mark_end(lexer);
            advance(lexer);
            if (lexer->lookahead == ']') {
                advance(lexer);
                if (lexer->lookahead == '>') {
                    advance(lexer);
                    if (advanced_once) {
                        lexer->result_symbol = CHAR_DATA;
                        return false;
                    }
                }
            }
        }
        advanced_once = true;
        advance(lexer);
    }

    if (advanced_once) {
        lexer->mark_end(lexer);
        lexer->result_symbol = CHAR_DATA;
        return true;
    }
    return false;
}

/// Scan for a CData node
static bool scan_cdata(TSLexer *lexer) {
    bool advanced_once = false;

    while (!lexer->eof(lexer)) {
        if (lexer->lookahead == ']') {
            lexer->mark_end(lexer);
            advance(lexer);
            if (lexer->lookahead == ']') {
                advance(lexer);
                if (lexer->lookahead == '>' && advanced_once) {
                    lexer->result_symbol = CDATA;
                    return true;
                }
            }
        }
        advanced_once = true;
        advance(lexer);
    }

    return false;
}

bool tree_sitter_xml_external_scanner_scan(void *payload, TSLexer *lexer, const bool *valid_symbols) {
    Vector *tags = (Vector *)payload;

    if (in_error_recovery(valid_symbols)) {
        return false;
    }

    if (valid_symbols[PI_TARGET]) {
        return scan_pi_target(lexer, valid_symbols);
    }

    if (valid_symbols[PI_CONTENT]) {
        return scan_pi_content(lexer);
    }

    if (valid_symbols[CHAR_DATA] && scan_char_data(lexer)) {
        return true;
    }

    if (valid_symbols[CDATA] && scan_cdata(lexer)) {
        return true;
    }

    switch (lexer->lookahead) {
        case '<':
            lexer->mark_end(lexer);
            advance(lexer);
            if (lexer->lookahead == '!') {
                advance(lexer);
                return scan_comment(lexer);
            }
            break;
        case '/':
            if (valid_symbols[SELF_CLOSING_TAG_DELIMITER]) {
                return scan_self_closing_tag_delimiter(tags, lexer);
            }
            break;
        case '\0':
            break;
        default:
            if (valid_symbols[START_TAG_NAME]) {
                return scan_start_tag_name(tags, lexer);
            }
            if (valid_symbols[END_TAG_NAME]) {
                return scan_end_tag_name(tags, lexer);
            }
    }

    return false;
}

void *tree_sitter_xml_external_scanner_create() {
    Vector *tags = (Vector *)ts_calloc(1, sizeof(Vector));
    if (tags == NULL) abort();
    array_init(tags);
    return tags;
}

void tree_sitter_xml_external_scanner_destroy(void *payload) {
    Vector *tags = (Vector *)payload;
    for (uint32_t i = 0; i < tags->size; ++i) {
        array_delete(array_get(tags, i));
    }
    array_delete(tags);
    ts_free(tags);
}

unsigned tree_sitter_xml_external_scanner_serialize(void *payload, char *buffer) {
    Vector *tags = (Vector *)payload;
    uint32_t tag_count = tags->size > UINT16_MAX ? UINT16_MAX : tags->size;
    uint32_t serialized_tag_count = 0, size = sizeof tag_count;

    memcpy(&buffer[size], &tag_count, size);
    size += sizeof tag_count;

    for (; serialized_tag_count < tag_count; ++serialized_tag_count) {
        String *tag = array_get(tags, serialized_tag_count);
        uint32_t name_length = tag->size;
        if (name_length > UINT8_MAX) {
            name_length = UINT8_MAX;
        }
        if (size + 2 + name_length >= TREE_SITTER_SERIALIZATION_BUFFER_SIZE) {
            break;
        }
        buffer[size++] = (char)name_length;
        if (name_length > 0) {
            memcpy(&buffer[size], tag->contents, name_length);
        }
        array_delete(tag);
        size += name_length;
    }

    memcpy(&buffer[0], &serialized_tag_count, sizeof serialized_tag_count);
    return size;
}

void tree_sitter_xml_external_scanner_deserialize(void *payload, const char *buffer, unsigned length) {
    Vector *tags = (Vector *)payload;

    for (unsigned i = 0; i < tags->size; ++i) {
        array_delete(array_get(tags, i));
    }
    array_delete(tags);

    if (length == 0) return;

    uint32_t size = 0, tag_count = 0, serialized_tag_count = 0;
    memcpy(&serialized_tag_count, &buffer[size], sizeof serialized_tag_count);
    size += sizeof serialized_tag_count;
    memcpy(&tag_count, &buffer[size], sizeof tag_count);
    size += sizeof tag_count;

    if (tag_count == 0) return;

    array_reserve(tags, tag_count);

    uint32_t iter = 0;
    for (; iter < serialized_tag_count; ++iter) {
        String tag = array_new();
        tag.size = (uint8_t)buffer[size++];
        if (tag.size > 0) {
            array_reserve(&tag, tag.size + 1);
            memcpy(tag.contents, &buffer[size], tag.size);
            size += tag.size;
        }
        array_push(tags, tag);
    }
    // add zero tags if we didn't read enough, this is because the
    // buffer had no more room but we held more tags.
    for (; iter < tag_count; ++iter) {
        String tag = array_new();
        array_push(tags, tag);
    }
}
