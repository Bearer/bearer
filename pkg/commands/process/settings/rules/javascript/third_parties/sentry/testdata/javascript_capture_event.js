Sentry.captureEvent({
  message: "user successfully logged in " + current_user.email,
  stacktrace: [
    // ...
  ],
})
