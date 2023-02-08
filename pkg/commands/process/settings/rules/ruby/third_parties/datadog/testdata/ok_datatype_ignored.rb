Datadog.configure do |c|
  user = { user_id: 42 }
  c.tags = user
end

Datadog::Tracing.trace("web.request", tags: { user_id: 42 }) do |span, trace|
  call
end
