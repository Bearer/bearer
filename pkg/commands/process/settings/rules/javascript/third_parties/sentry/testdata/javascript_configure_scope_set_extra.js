Sentry.configureScope((scope) => {
	scope.setExtra("email", user.email);
});
