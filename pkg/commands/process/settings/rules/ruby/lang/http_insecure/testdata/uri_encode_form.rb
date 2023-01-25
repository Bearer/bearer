uri = URI('http://my.api.com/users/search')
user = { email: current_user.email }
uri.query = URI.encode_www_form(user)