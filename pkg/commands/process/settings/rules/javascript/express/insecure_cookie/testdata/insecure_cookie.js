const session = require("express-session");
const express = require("express");
const app = express();

app.use(
	session({
		cookie: {
			domain: "example.com",
			secure: false, // Ensures the browser only sends the cookie over HTTPS.
			httpOnly: false,
			name: "my-custom-cookie-name",
			maxAge: 24 * 60 * 60 * 1000,
			path: "/some-path"
		},
	})
);