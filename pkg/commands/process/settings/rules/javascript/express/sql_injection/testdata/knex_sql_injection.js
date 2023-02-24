const knex = require("knex")({
	client: "mysql",
});

module.exports.badQuery = function (req, res) {
	var cartDetails = knex
		.select("user.cart_details")
		.from("users")
		.whereRaw("name = " + req.query.user.name);

	res.send(prepareJson(cartDetails));
};
