cipher = OpenSSL::Cipher.new('aes-128-cbc')
dsa_encrypt = OpenSSL::PKey::DSA.new(2048)
dsa_encrypt.export(cipher, "hello world")

OpenSSL::PKey::DSA.new(2048).to_pem(cipher, "hello world")