url = "http://#{ ENV['ORDER_SERVICE_DOMAIN'] }/whatever"
url = "http://" + domain + ".com" + "/other"
order_service_url = ENV["ORDER_SERVICE_URL"]
user_service_host = ENV.fetch('USER_SERVICE_HOST', "default")
account_id = ENV["ACCOUNT_ID"]
other = other["IGNORE_ME_HOST"]

x = { "ignore.domain.com" => "abc" }
y = x["ignore.domain.com"]
