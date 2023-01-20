# trigger_condition: processing sensitive data
user.gender_identity

# trigger:verification disabled
http.verify_mode = OpenSSL::SSL::VERIFY_NONE

uri = URI('https://secure.example.com/some_path?query=string')
Net::HTTP.start(uri.host, uri.port, :use_ssl => true, :verify_mode => OpenSSL::SSL::VERIFY_NONE) do |http|
  Net::HTTP::Get.new uri
end

# ok:verification enabled
http.verify_mode = OpenSSL::SSL::VERIFY_PEER
