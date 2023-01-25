uri = URI("https://my.api.com/users/search")

response = Curl.get("https://my.api.com/users/search")

RestClient.get("https://my.api.com/users/search")

Typhoeus.get("https://my.api.com/users/search")

HTTParty.get("https://my.api.com/users/search")

HTTP.get("https://my.api.com/users/search", params: { user: { id: 12 } })

Excon.get("https://my.api.com/users/search")

Faraday.get("https://my.api.com/users/search")

HTTPX.get("https://my.api.com/users/search")