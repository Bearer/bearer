const session = require("cookie-session");
const express = require("express");
const app = express();

app.use(
	session({
		cookie: {
			domain: "example.com",
			secure: true,
			httpOnly: false,
			maxAge: 24 * 60 * 60 * 1000,
			path: "/some-path",
      name: "my-custom-cookie-name"
		},
	})
);
