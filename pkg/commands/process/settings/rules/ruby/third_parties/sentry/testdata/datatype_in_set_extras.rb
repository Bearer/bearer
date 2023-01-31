Sentry.set_extras(email: user.email)

Sentry.configure_scope do |scope|
  scope.set_extras(email: user.email)
end

Sentry.with_scope do |scope|
  scope.set_extras(email: user.email)
end
