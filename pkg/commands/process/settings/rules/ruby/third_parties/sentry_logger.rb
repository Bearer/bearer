# trigger_condition: application has sensitive data
class User
  attr_accessor :email
end

current_user = user = User.new

# trigger: logging data type to Sentry
Sentry::Breadcrumb.new(
  category: "auth",
  message: "Authenticated user #{user.email}",
  level: "info"
)

# trigger: data type in Sentry breadcrumb init
Sentry.init do |config|
  config.before_breadcrumb = lambda do |breadcrumb, hint|
    breadcrumb.message = "Authenticated user #{current_user.email}"
    breadcrumb
  end
end

# trigger: sending data type to Sentry.set_user
# https://docs.sentry.io/platforms/ruby/guides/rails/enriching-events/identify-user/
Sentry.set_user(email: user.email)
