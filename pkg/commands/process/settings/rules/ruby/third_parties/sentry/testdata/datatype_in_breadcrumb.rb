Sentry::Breadcrumb.new(
  category: "auth",
  message: "Authenticated user #{user.email}",
  level: "info"
)
