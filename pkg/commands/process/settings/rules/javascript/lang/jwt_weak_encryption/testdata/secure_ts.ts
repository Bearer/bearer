var jwt = require("jsonwebtoken");

var token = jwt.sign({ foo: "bar" }, process.env.JWT_SECRET, {
	algorithm: "ES256",
});
