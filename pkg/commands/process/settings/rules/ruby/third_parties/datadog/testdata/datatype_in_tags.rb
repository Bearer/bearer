Datadog.configure do |c|
  user = { email: "someone@example.com" }
  c.tags = user
end

span = Datadog.configuration[:cucumber][:tracer].active_span
span.set_tag('user.email', user.email)

Datadog::Tracing.active_span&.set_tag('customer.id', user.email)
Datadog::Tracing.active_span.set_tag('customer.id', user.email)

Datadog::Tracing.trace("web.request", tags: { email: user.email }) do |span, trace|
  call
end

Datadog::Tracing.trace("web.request") do |span, trace|
  span.set_tag('user.email', user.email)
end
