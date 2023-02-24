const knex = require("knex")({
	client: "sqlite3",
	connection: {
		filename: "./data.db",
	},
});

app.post("/users", (req, res) => {
	knex("users").select(knex.raw("count " + req.query.user_field + ""));
});
