begin
  raise CustomException.new(current_user.email)
rescue => exception
  Bugsnag.notify(exception)
end

Bugsnag.notify(exception) do |event|
  # Adjust the severity of this error
  event.severity = "warning"

  # Add customer information to this event
  event.add_metadata(:account, {
    user_name: current_user.name,
    paying_customer: true
  })
end