# Ruby third-party data send

class User
  attr_accessor :email
end

current_user = user = User.new

## Detected
Sentry::Breadcrumb.new(
  category: "auth",
  message: "Authenticated user #{user.email}",
  level: "info"
)

Sentry.init do |config|
  config.before_breadcrumb = lambda do |breadcrumb, hint|
    breadcrumb.message = "Authenticated user #{current_user.email}"
    breadcrumb
  end
end

# https://docs.sentry.io/platforms/ruby/guides/rails/enriching-events/identify-user/
Sentry.set_user(email: user.email)
