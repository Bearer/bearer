var jwt = require("jsonwebtoken");

var token = jwt.sign(
	{ user: { uuid: "1fbae5ff-86c8-4ece-8278-bd94957de1bf" } },
	process.env.JWT_SECRET
);
