Datadog.configure do |c|
  user = { foo: "bar" }
  c.tags = user
end

Datadog::Tracing.trace("web.request", tags: { foo: "bar" }) do |span, trace|
  call
end
