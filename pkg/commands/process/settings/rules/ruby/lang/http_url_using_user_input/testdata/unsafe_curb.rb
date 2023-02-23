# https://www.rubydoc.info/github/taf2/curb/

Curl.http("GET", params[:oops], nil) do
end

Curl.get(params[:oops]) {}

Curl::Easy.perform(params[:oops]) {}

easy = Curl::Easy.new(params[:oops]) {}
easy.url = params[:oops2]

easy2 = Curl::Easy.new
easy2.url = params[:oops]

Curl::Multi.get(["https://my.api.com/secure", params[:oops]], {}) {}

Curl::Multi.http([{ url: params[:oops], method: :post }]) {}
