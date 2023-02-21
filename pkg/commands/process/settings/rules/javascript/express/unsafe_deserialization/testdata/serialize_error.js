import deserializeError from 'serialize-error';

module.exports.deserializedError = function(req, _res) {
  deserializeError({
    name: "MyCustomError",
    message: req.params.error
  })
}