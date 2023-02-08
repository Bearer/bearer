# rule ignores Unique Identifier datatypes
Airbrake.notify('App crashed') do |notice|
  notice[:params][:email] = customer.id
end