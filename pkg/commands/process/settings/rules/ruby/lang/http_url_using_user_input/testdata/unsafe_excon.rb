connection = Excon.new(params[:oops], foo: true)
connection2 = Excon.new("http://example.com", path: params[:oops])

connection3 = Excon::Connection.new(host: params[:oops])
connection4 = Excon::Connection.new(hostname: params[:oops])
connection5 = Excon::Connection.new(path: params[:oops])
connection6 = Excon::Connection.new(port: params[:oops])

connection.post(path: params[:oops])

connection2.request(path: params[:oops])

connection3.requests([{ :method => :get, path: params[:oops] }])

Excon.get(params[:oops])
Excon.post("http://example.com", path: params[:oops])
