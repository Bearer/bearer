# add attributes
def track_user(user)
  current_span = OpenTelemetry::Trace.current_span

  current_span.add_attributes({
    "user.id" => user.id,
    "user.first_name" => user.first_name
  })
end

# add attributes at span creation
Tracer.in_span("data leaking", attributes: { "current_user" => user.email, "date" => DateTime.now }) do |span|
  puts "in the span block"
end

# set attributes
current_span = OpenTelemetry::Trace.current_span
users = [user.email, customer.email, admin_user.email]
current_span.set_attribute("current_users", users)


