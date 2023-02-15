const StatsD = require("hot-shots");
const client = new StatsD({
	port: 8020,
	globalTags: { env: process.env.NODE_ENV },
	errorHandler: errorHandler,
});

client.event("user", "logged_in", {});
