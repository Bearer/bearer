import { Sequelize } from "sequelize";

module.exports.fooBar = function (req, _res) {
	var sqlite = new Sequelize("sqlite::memory:");
	var customerQuery =
		"SELECT * FROM customers WHERE status = " + req.params.customer.status;
	sqlite.query(customerQuery);
};
