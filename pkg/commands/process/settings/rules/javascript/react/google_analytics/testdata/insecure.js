const user = { email: "jhon@gmail.com " };
ReactGA.event({
	category: "user",
	action: "logged_in",
	value: user.email,
});
