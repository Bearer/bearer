# https://docs.ruby-lang.org/en/master/Net/HTTP.html

response = Net::HTTP.post_form("http://my.api.com/users/search")

Net::HTTP.start("http://my.api.com/users/search") do
end

Net::HTTP::Get.new("http://my.api.com/users/search", { "X-Test": 42 })
