const vm = require('node:vm');

exports.handler = async function(event, _context) {
  const context = event["params"]["context"];
  vm.createContext("count");

  var ok = new vm.Script("count += 1")
}
