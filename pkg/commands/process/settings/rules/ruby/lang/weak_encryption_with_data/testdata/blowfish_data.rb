Crypt::Blowfish.new("insecure").encrypt_block do |user|
  user.password
end