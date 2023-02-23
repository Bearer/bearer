const { Client } = require('pg')
const pgClient = new Client({
  // pg client setup
})

const knex = require('knex')({
  client: 'mysql',
})

const connection = mysql.createConnection({});
const asyncConn = await mysql.createConnection({});

exports.handler =  async function(_event, _context) {
  var user = getCurrentUser()
  var userRes = pgClient.query('SELECT * FROM users WHERE user.name = ' + user.name)
  var res = knex.select('user.cart_details')
    .from('users')
    .whereRaw('id = '+ userRes.id)

  connection.query("SELECT * FROM `user` WHERE name = " + userRes.name);

  return res
}
