uri = URI('http://my.api.com/users/search')
user = { x: 42 }
uri.query = URI.encode_www_form(user)

response = Net::HTTP.post_form(uri, { x: "abc" })

uri = URI('https://my.api.com/users/search')
