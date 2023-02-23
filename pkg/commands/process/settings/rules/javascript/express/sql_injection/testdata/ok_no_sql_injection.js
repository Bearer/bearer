import { Sequelize } from "sequelize";
const { Client } = require('pg')
const client = new Client({
  // pg client setup
})

const connection = mysql.createConnection({});

module.exports.fooBar = function(req, _res) {
  var sqlite = new Sequelize('sqlite::memory:')
  var customerQuery = "SELECT * FROM customers WHERE status = ACTIVE"
  sqlite.query(customerQuery)

  client.query('SELECT * FROM users WHERE user.name = ' + getUser().name)

  connection.query("SELECT * FROM `user` WHERE name = " + currentUser().name);
}
