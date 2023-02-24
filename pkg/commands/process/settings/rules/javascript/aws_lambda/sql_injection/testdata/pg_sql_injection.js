const { Client } = require('pg')

const client = new Client({
  // client setup
})

exports.handler =  async function(event, _context) {
  var user = client.query('SELECT * FROM users WHERE user.name = ' + event.user.name)

  return user
}