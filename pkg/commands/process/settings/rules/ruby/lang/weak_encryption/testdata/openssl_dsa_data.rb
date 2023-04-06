cipher = OpenSSL::Cipher.new('aes-128-cbc')
dsa_encrypt = OpenSSL::PKey::DSA.new(2048)
dsa_encrypt.export(cipher, user.email)

OpenSSL::PKey::DSA.new(2048).to_pem(cipher, user.first_name)