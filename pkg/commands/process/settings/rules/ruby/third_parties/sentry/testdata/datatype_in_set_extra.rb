Sentry.configure_scope do |scope|
  scope.set_extra(:email, user.email)
end

Sentry.with_scope do |scope|
  scope.set_extra(:email, user.email)
end
