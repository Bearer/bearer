import myJWT from "jsonwebtoken";

import {jwt as myJWT} from "jsonwebtoken";

const myJWT = require("jsonwebtoken").jwt;

const privateKey = "foo";
myJWT.sign(user, privateKey, {
	expiresInMinutes: 60 * 5,
	algorithm: "RS256",
});
