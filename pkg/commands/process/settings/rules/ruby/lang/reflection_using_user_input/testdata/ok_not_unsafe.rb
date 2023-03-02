Object.const_get(x)
Object.const_set(x, 42)
Object.remove_const(x)

method(m)

x.define_method(m) {}

m.to_sym.to_proc

bad_things(&m.to_sym)
x.bad_things(&m.to_sym)


x.constantize
