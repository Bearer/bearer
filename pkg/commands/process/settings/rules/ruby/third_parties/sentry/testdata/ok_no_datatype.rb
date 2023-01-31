Sentry::Breadcrumb.new(
  category: "tracker",
  message: "some event has happened",
  level: "info"
)

Sentry.capture_message("test")

Sentry.init do |config|
  config.before_breadcrumb = lambda do |breadcrumb, hint|
    breadcrumb
  end
end

Sentry.set_context(something: "hey")
Sentry.set_extras(something: "hey")
Sentry.set_tags(something: "hey")
Sentry.set_user(something: "hey")

Sentry.configure_scope do |scope|
  scope.set_context('something', { hey: "there" })
  scope.set_extra('something', "hey")
  scope.set_extras(something: "hey")
  scope.set_tag('something', "hey")
  scope.set_tags(something: "hey")
  scope.set_user(something: "hey")
end

Sentry.with_scope do |scope|
  scope.set_context('something', { hey: "there" })
  scope.set_extra('something', "hey")
  scope.set_extras(something: "hey")
  scope.set_tag('something', "hey")
  scope.set_tags(something: "hey")
  scope.set_user(something: "hey")
end
