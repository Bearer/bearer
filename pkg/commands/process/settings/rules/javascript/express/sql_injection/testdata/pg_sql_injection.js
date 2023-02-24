const { Client } = require("pg");

const client = new Client({
	// client setup
});

module.exports.fooBar = function (req, _res) {
	var user = client.query(
		"SELECT * FROM users WHERE user.name = " + req.params.user.name
	);

	return user;
};
