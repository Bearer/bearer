Digest::SHA1.hexidigest(user.password)
Digest::MD5.hexdigest(user.password)

RC4.new("insecure").encrypt(user.password)
Crypt::Blowfish.new("insecure").encrypt_block({ |u| user.password })

OpenSSL::PKey::RSA.new(File.read('rsa.pem')).private_encrypt(customer.password)

cipher = OpenSSL::Cipher.new('aes-128-cbc')
OpenSSL::PKey::DSA.new(2048).export(cipher, customer.password)
OpenSSL::PKey::RSA.new(2048).to_pem(cipher, customer.password)

# Not detected (not a password)
Digest::SHA1.hexidigest(user.email)
Digest::MD5.hexdigest(user.first_name)
RC4.new("insecure").encrypt(user.address)
Crypt::Blowfish.new("insecure").encrypt_block({ |u| user.gender_identity })

# Not detected (variable assignment)
pkey = OpenSSL::PKey::RSA.new(File.read('rsa.pem'))
pkey.private_encrypt(customer.password)