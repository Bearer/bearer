RC4.new("insecure").encrypt(user.password)

rc4_encrypt = RC4.new("insecure")
rc4_encrypt.encrypt!(user.password)

