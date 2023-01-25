OpenSSL::PKey::RSA.new(File.read('rsa.pem')).private_encrypt("test")

cipher = OpenSSL::Cipher.new('aes-128-cbc')
rsa_encrypt = OpenSSL::PKey::RSA.new(2048)
rsa_encrypt.export(cipher, "hello world")

OpenSSL::PKey::RSA.new(2048).to_pem(cipher, "hello world")