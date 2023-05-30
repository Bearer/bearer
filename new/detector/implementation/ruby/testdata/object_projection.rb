let obj = { x: { a: { i: 3 } }, y: 4 }

# Known properties
obj.x
obj["x"].a

# Unknown properties
obj.z
@myvar.x
@myvar["w"]

# Multiple index
foo = [:a, :b, :c, :d, :e]
foo[0, 2]