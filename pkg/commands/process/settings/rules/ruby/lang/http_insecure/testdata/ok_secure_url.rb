response = Net::HTTP.post_form("https://my.api.com/users/search")


Curl.http("GET", "https://my.api.com/users/search", nil) do
end

Curl::Multi.http([
  {
    url: "https://my.api.com/users/search",
    post_fields: { x: "http://my.api.com/users/search" }
  }
]) {}


connection = Excon.new("https://my.api.com/secure", foo: true)

Excon.get("https://my.api.com/secure", foo: true)


Faraday.get("https://api.secure.com")
Faraday.post("https://api.secure.com")
