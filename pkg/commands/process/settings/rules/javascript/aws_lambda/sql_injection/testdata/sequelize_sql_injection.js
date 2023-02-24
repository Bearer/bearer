import { Sequelize } from "sequelize";

exports.handler =  async function(event, _context) {
  var sqlite = new Sequelize('sqlite::memory:')
  var customerQuery = "SELECT * FROM customers WHERE status = " + event.customer.status
  sqlite.query(customerQuery)
}