uri = URI("http://my.api.com/users/search?ethnic_origin=#{user_1.ethnic_origin}")

uri = URI('http://my.api.com/users/search')
user = { first_name: "John", last_name: "Doe" }
uri.query = URI.encode_www_form(user)

response = Net::HTTP.post_form(uri, { user: { first_name: "John", last_name: "Doe" } })
