Sentry.configureScope((scope) => {
  scope.setTag("user_email", user.email)
})
