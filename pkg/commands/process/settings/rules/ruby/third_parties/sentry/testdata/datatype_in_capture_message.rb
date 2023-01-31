Sentry.capture_message("test: #{user.email}")
Sentry.capture_message("test", extra: { email: user.email })
Sentry.capture_message("test", tags: { email: user.email })
Sentry.capture_message("test", user: { email: user.email })
