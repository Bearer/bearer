const config = {
	clientID: process.env["GOOGLE_CLIENT_ID"],
	clientSecret: "secretHardcodedString",
	callbackURL: "/oauth2/redirect/google",
	scope: ["profile"],
};
