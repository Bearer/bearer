response = Net::HTTP.post_form("http://my.api.com/users/search", email: user.email)
