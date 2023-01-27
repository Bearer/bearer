Crypt::Blowfish.new("insecure").encrypt_block do |user|
  user.password
end

Crypt::Blowfish.new("your-key").encrypt_string(user.email)