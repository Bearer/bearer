Sentry.configureScope((scope) => {
  scope.setUser({ email: user.email })
})
