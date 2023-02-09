# add attributes at span creation
Tracer.in_span("data leaking", attributes: { "current_user" => user.email, "date" => DateTime.now }) do |span|
  puts "in the span block"
end

SomeOtherTracer.in_span("data leaking", attributes: { "current_user" => user.email, "date" => DateTime.now }) do |span|
  span.add_attributes(user.email)
end

YetAnotherTracer.in_span("data leaking", attributes: { "date" => DateTime.now }) do |span|
  span.add_event("leaking data for #{user.email}")
end
