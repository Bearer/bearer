span = OpenTelemetry::Trace.current_span
span.add_event("Schedule job for user: #{current_user.email}")
span.add_event("Cancel job for user", attributes: {
  "current_user" => current_user.email
})