let obj = { x: { a: { i: 3 } }, y: 4 }

// Known properties
obj.x
obj["x"].a

// Unknown properties
obj.z
obj["w"]

// Call
obj.x(i).a
