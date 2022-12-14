# Detected
Digest::SHA1.hexidigest(user.email)
Digest::MD5.hexdigest(user.first_name)
RC4.new("insecure").encrypt(user.address)

Crypt::Blowfish.new("insecure").encrypt_block({ |u| user.gender_identity })
