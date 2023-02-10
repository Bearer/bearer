var jwt = require("jsonwebtoken");
var token = jwt.sign({ user: { email: "jhon@gmail.com" } }, "shhhhh");
