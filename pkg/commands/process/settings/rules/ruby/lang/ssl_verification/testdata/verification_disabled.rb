http.verify_mode = OpenSSL::SSL::VERIFY_NONE

uri = URI('https://secure.example.com/some_path?query=string')
Net::HTTP.start(uri.host, uri.port, :use_ssl => true, :verify_mode => OpenSSL::SSL::VERIFY_NONE) do |http|
  Net::HTTP::Get.new uri
end