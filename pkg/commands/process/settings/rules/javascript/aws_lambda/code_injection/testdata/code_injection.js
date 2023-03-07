const vm = require('node:vm');

exports.handler = async function(event, _context) {
  const context = event["params"]["context"];
  vm.createContext(context);

  var bad = new vm.Script(event["query"])
}
