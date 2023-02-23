const { Client } = require('pg')
const pgClient = new Client({
  // pg client setup
})

const knex = require('knex')({
  client: 'mysql',
})

exports.handler =  async function(_event, _context) {
  var user = getCurrentUser()
  var userRes = pgClient.query('SELECT * FROM users WHERE user.name = ' + user.name)
  var res = knex.select('user.cart_details')
    .from('users')
    .whereRaw('id = '+ userRes.id)

  return res
}
