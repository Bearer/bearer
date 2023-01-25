OpenSSL::PKey::RSA.new(File.read('rsa.pem')).private_encrypt(user.email)

cipher = OpenSSL::Cipher.new('aes-128-cbc')
rsa_encrypt = OpenSSL::PKey::RSA.new(2048)
rsa_encrypt.export(cipher, user.email)

OpenSSL::PKey::RSA.new(2048).to_pem(cipher, user.email)