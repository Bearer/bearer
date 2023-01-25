blowfish = Crypt::Blowfish.new("insecure")
blowfish.encrypt_block do
  "hello world"
end

