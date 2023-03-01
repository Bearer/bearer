const session = require("express-session");
const express = require("express");
const app = express();

app.use(
	session({
    name: "my-custom-session-name",
		cookie: {
			domain: "example.com",
			secure: true,
			httpOnly: false,
			maxAge: 24 * 60 * 60 * 1000,
			path: "/some-path"
		},
	})
);
