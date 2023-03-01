const session = require("express-session");
const express = require("express");
const app = express();

app.use(
	session({ name: "my-custom-session-name" })
);
