Sentry::Breadcrumb.new(
  category: "auth",
  message: "user has authenticated #{current_user.id}",
  level: "info"
)
