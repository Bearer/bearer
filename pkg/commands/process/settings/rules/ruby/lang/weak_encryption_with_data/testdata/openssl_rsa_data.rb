OpenSSL::PKey::RSA.new(File.read('rsa.pem')).private_encrypt(user.password)

cipher = OpenSSL::Cipher.new('aes-128-cbc')
rsa_encrypt = OpenSSL::PKey::RSA.new(2048)
rsa_encrypt.export(cipher, user.password)

OpenSSL::PKey::RSA.new(2048).to_pem(cipher, user.first_name)