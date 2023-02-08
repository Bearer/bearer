begin
  1/0
rescue ZeroDivisionError => ex
  response = Airbrake.notify_sync(ex)
end
