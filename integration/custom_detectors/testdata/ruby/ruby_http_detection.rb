## URI

uri = URI("http://my.api.com/users/search?ethnic_origin=#{user_1.ethnic_origin}")

uri = URI('https://my.api.com/users/search')
user_1 = { first_name: "John", last_name: "Doe" }
uri.query = URI.encode_www_form(user_1)


## Net::HTTP

response = Net::HTTP.post_form(uri, { user_2: { first_name: "John", last_name: "Doe" } })


## Curl

User = Struct.new(:first_name, :last_name, keyword_init: true)
user_3 = User.new(first_name: "first", last_name: "last")
response = Curl.get("http://my.api.com/users/search?first_name=#{user_2.first_name}")

response = Curl.get("https://my.api.com/users/search?first_name=#{user_2.first_name}")

user_4 = { first_name: "John", last_name: "Doe" }
response = Curl.post("http://my.api.com/users/create", user_4)

response = Curl.post("http://my.api.com/users/create", { user_5: { first_name: "John", last_name: "Doe" } })


## RestClient

RestClient.get("http://my.api.com/users/search?first_name=#{user_2.first_name}")

RestClient.get("https://my.api.com/users/search?first_name=#{user_2.first_name}")

RestClient.post("http://my.api.com/users/create", { user_6: { first_name: "John", last_name: "Doe" } })


## Typhoeus

options = { body: { user_7: { first_name: "John", last_name: "Doe" } } }
response = Typhoeus.post("http://my.api.com/users/create", options)

response = Typhoeus.post("http://my.api.com/users/create", { body: { user_8: { first_name: "John", last_name: "Doe" } } })

Typhoeus.get("http://my.api.com/users/search?first_name=#{user_9.first_name}")

Typhoeus.get("https://my.api.com/users/search?first_name=#{user_9.first_name}")


## HTTParty

HTTParty.get("http://my.api.com/users/search?first_name=#{user_2.first_name}")

HTTParty.get("https://my.api.com/users/search?first_name=#{user_2.first_name}")

params = {
	body: {
    user: {
      first_name: "John",
      last_name: "Doe",
    }
	}
}
HTTParty.post("http://my.api.com/users/create", params)

HTTParty.post("http://my.api.com/users/create", { body: { user: { first_name: "John", last_name: "Doe" } } })


## HTTP.rb

HTTP.get("http://my.api.com/users/search?first_name=#{user_2.first_name}")

HTTP.get("https://my.api.com/users/search", params: { user_8: { first_name: "John" } })

HTTP.post("http://my.api.com/users/create", form: { user_9: { first_name: "John", last_name: "Doe" } })


## Excon

Excon.get("http://my.api.com/users/search?first_name=#{user_2.first_name}")

Excon.get("https://my.api.com/users/search?first_name=#{user_2.first_name}")

Excon.post("http://my.api.com/users/create", body: { user_10: { first_name: "John", last_name: "Doe" } })


## Faraday

Faraday.get("http://my.api.com/users/search?first_name=#{user_2.first_name}")

Faraday.get("https://my.api.com/users/search?first_name=#{user_2.first_name}")

params_2 = { user_11: { first_name: "John", last_name: "Doe" } }

encoded_params = URI.encode_www_form(params_2)

response = Faraday.post("http://my.api.com/users/create", encoded_params)

# Faraday.post("http://my.api.com/users/create") do |request|
#   request.body = { user: { first_name: "John", last_name: "Doe" } }.to_json
# end


## HTTPX

HTTPX.post("http://my.api.com/users/create", json: { user_12: { first_name: "John", last_name: "Doe" } })

HTTPX.get("http://my.api.com/users/search?first_name=#{user_2.first_name}")

HTTPX.get("https://my.api.com/users/search?first_name=#{user_2.first_name}")
