let obj = { x: { a: { i: 3 } }, y: 4 }

# Known properties
obj.x
obj["x"].a

# Unknown properties
obj.z
@myvar.x
@myvar["w"]

# Call with arguments
obj.x({ email: " " }, { first_name: "" })
obj.x({ email: " " }, { first_name: "" }).a
