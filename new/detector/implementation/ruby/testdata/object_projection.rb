let obj = { x: { a: { i: 3 } }, y: 4 }

# Known properties
obj.x
obj["x"].a

# Unknown properties
obj.z
@myvar.x
@myvar["w"]
