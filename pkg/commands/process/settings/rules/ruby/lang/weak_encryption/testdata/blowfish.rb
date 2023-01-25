blowfish = Crypt::Blowfish.new("insecure")
blowfish.encrypt_block({ |u| user.gender_identity })

