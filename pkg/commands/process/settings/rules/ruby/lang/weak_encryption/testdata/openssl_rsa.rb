OpenSSL::PKey::RSA.new(File.read('rsa.pem')).private_encrypt("test")

size = 2048
cipher = OpenSSL::Cipher.new('aes-128-cbc')
rsa_encrypt = OpenSSL::PKey::RSA.new(size)
rsa_encrypt.export(cipher, "hello world") # secure

low_size = 512
OpenSSL::PKey::RSA.new(2048).to_pem(cipher, "hello world") # secure
OpenSSL::PKey::RSA.new(512)
OpenSSL::PKey::RSA.new(low_size)
OpenSSL::PKey::RSA.new(512).to_pem(cipher, "hello world") # unsecure
OpenSSL::PKey::RSA.new(low_size).to_pem(cipher, "hello world") #unsecure