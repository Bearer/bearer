size = 2048
cipher = OpenSSL::Cipher.new('aes-128-cbc')
rsa_encrypt = OpenSSL::PKey::RSA.new(size)
dsa_encrypt.export(cipher, "hello world")

low_size = 512
OpenSSL::PKey::RSA.new(2048).to_pem(cipher, "hello world")
OpenSSL::PKey::RSA.new(512)
OpenSSL::PKey::RSA.new(low_size)
OpenSSL::PKey::RSA.new(512).to_pem(cipher, "hello world")
OpenSSL::PKey::RSA.new(low_size).to_pem(cipher, "hello world")