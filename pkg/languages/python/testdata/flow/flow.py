def with_statement():
    with source() as value, other:
        cursor_sink(value)
    
def for_statement():
    for value in source():
        result_sink(value)
        cursor_sink(value) # no match

cursor_sink(value) # no match
