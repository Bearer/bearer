# https://www.rubydoc.info/github/taf2/curb/

Curl.http("GET", "http://my.api.com/users/search", nil) do
end

Curl.get("http://my.api.com/users/search") {}

Curl::Easy.perform("http://my.api.com/users/search") {}

easy = Curl::Easy.new("http://my.api.com/users/search") {}
easy.url = "http://my.api.com/customers"

easy2 = Curl::Easy.new
easy2.url = "http://my.api.com/users/search"

Curl::Multi.get(["https://my.api.com/secure", "http://my.api.com/users/search"], {}) {}

Curl::Multi.http([{ url: "http://my.api.com/users/search", method: :post }]) {}
