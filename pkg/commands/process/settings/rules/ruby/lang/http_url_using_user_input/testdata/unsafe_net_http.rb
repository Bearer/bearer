# https://docs.ruby-lang.org/en/master/Net/HTTP.html

response = Net::HTTP.post_form("http://#{params[:oops]}/users/search")

Net::HTTP.start(params[:host]) do |instance1|
  instance1.head(params[:path])
end

Net::HTTP::Get.new(params[:oops], { "X-Test": 42 })

instance2 = Net::HTTP.start(params[:oops])
instance2.ipaddr = request.env[:oops]
instance2.send_request("GET", params[:oops], nil)

instance3 = Net::HTTP.new(params[:oops])
instance3.patch(params[:path])
instance3.start do |instance4|
  instance4.post(request.env[:oops])
end
