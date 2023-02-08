begin
  1/0
rescue ZeroDivisionError => ex
  response = Airbrake.notify_sync(current_user.first_name)
end