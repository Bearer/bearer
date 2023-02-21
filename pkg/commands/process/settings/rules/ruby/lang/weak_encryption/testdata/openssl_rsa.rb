OpenSSL::PKey::RSA.new(File.read('rsa.pem')).private_encrypt("test")

cipher = OpenSSL::Cipher.new('aes-128-cbc')
rsa_encrypt = OpenSSL::PKey::RSA.new(2048) # secure
rsa_encrypt.export(cipher, "hello world") # secure

OpenSSL::PKey::RSA.new(2048).to_pem(cipher, "hello world") # secure
OpenSSL::PKey::RSA.new(512).to_pem(cipher, "hello world") # unsecure
