def with_statement():
    with source() as value, other:
        cursor_sink(value)

def for_statement():
    for value in source():
        result_sink(value)
        cursor_sink(value) # no match

def reflexive_methods():
    s = source()
    x = s.format("hello")
    result_sink(x)
    cursor_sink(x) # no match

def non_reflexive_methods():
    s = source()
    x = s.my_method("hello")
    result_sink(x) # no match
    cursor_sink(x) # no match

cursor_sink(value) # no match
