uri_1 = URI("http://my.api.com/users/search?ethnic_origin=#{user_1.ethnic_origin}")

uri_2 = URI('http://my.api.com/users/search')
user = { first_name: "John", last_name: "Doe" }
uri_2.query = URI.encode_www_form(user)

response = Net::HTTP.post_form(uri_2, { user: { first_name: "John", last_name: "Doe" } })

uri_3 = URI('https://my.api.com/users/search')

Net::HTTP.post_form('http://my.api.com/users/search', { user: { first_name: "John", last_name: "Doe" } })
Net::HTTP.post_form('https://my.api.com/users/search', { user: { first_name: "John", last_name: "Doe" } })
