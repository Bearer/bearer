Sentry.init do |config|
  config.before_breadcrumb = lambda do |breadcrumb, hint|
    breadcrumb.message = "Authenticated user #{current_user.email}"
    breadcrumb
  end
end
