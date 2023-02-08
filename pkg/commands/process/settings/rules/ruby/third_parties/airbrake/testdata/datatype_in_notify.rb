Airbrake.notify(user.first_name)

Airbrake.notify('App crashed!', {
  current_user: user.email
})

Airbrake.notify('App crashed') do |notice|
  notice[:params][:email] = customer.email
end