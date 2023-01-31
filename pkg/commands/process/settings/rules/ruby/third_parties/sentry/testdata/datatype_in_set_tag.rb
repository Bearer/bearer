Sentry.configure_scope do |scope|
  scope.set_tag(:email, user.email)
end

Sentry.with_scope do |scope|
  scope.set_tag(:email, user.email)
end
