# sending data type to Sentry.set_user
# https://docs.sentry.io/platforms/ruby/guides/rails/enriching-events/identify-user/
Sentry.set_user(email: user.email)

Sentry.configure_scope do |scope|
  scope.set_user(email: user.email)
end

Sentry.with_scope do |scope|
  scope.set_user(email: user.email)
end
