Curl.http("GET", x, nil) do
end

Curl.get(x) {}

Curl::Easy.perform(x) {}

easy = Curl::Easy.new(x) {}
easy.url = y

easy2 = Curl::Easy.new
easy2.url = x

Curl::Multi.get(["https://my.api.com/secure", x], {}) {}

Curl::Multi.http([{ url: x, method: :post }]) {}


connection = Excon.new(x, foo: true)
connection2 = Excon.new("http://example.com", path: x)

connection3 = Excon::Connection.new(host: x)
connection4 = Excon::Connection.new(hostname: x)
connection5 = Excon::Connection.new(path: x)
connection6 = Excon::Connection.new(port: x)

connection.post(path: x)

connection2.request(path: x)

connection3.requests([{ :method => :get, path: x }])

Excon.get(x)
Excon.post("http://example.com", path: x)


Faraday.get(x)


response = Net::HTTP.post_form("http://#{x}/users/search")

Net::HTTP.start(x) do |instance1|
  instance1.head(x)
end

Net::HTTP::Get.new(x, { "X-Test": 42 })

instance2 = Net::HTTP.start(x)
instance2.ipaddr = x
instance2.send_request("GET", x, nil)

instance3 = Net::HTTP.new(x)
instance3.patch(x)
instance3.start do |instance4|
  instance4.post(x)
end


Faraday.post(x)
