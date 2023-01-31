Sentry.set_context('email', { email: user.email })

Sentry.configure_scope do |scope|
  scope.set_context('email', { email: user.email })
end

Sentry.with_scope do |scope|
  scope.set_context('email', { email: user.email })
end
