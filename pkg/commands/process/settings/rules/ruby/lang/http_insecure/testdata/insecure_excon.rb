connection = Excon.new("http://my.api.com/insecure", foo: true)

Excon.get("http://my.api.com/insecure", foo: true)
