cipher = OpenSSL::Cipher.new('aes-128-cbc')
rsa_encrypt = OpenSSL::PKey::RSA.new(2048)
rsa_encrypt.export(cipher, customer.password)

OpenSSL::PKey::RSA.new(2048).to_pem(cipher, customer.password)