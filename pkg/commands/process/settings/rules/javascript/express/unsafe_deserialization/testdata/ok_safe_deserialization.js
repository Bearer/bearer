import deserializeError from 'serialize-error';
var nodeSerialize = require("node-serialize")

module.exports.safeDeserialization = function(req, _res) {
  deserializeError({
    name: "MyCustomError",
    message: "Something went wrong"
  })

  nodeSerialize.unserialize({ hello: "world" })
}
