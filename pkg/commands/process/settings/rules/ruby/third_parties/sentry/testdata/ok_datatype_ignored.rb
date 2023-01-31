Sentry::Breadcrumb.new(
  category: "auth",
  message: "user has authenticated #{current_user.user_id}",
  level: "info"
)
