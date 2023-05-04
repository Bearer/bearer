class User
  def private_key
    @private_key ||= '-----BEGIN PGP PRIVATE KEY BLOCK-----asdf-----END PGP PRIVATE KEY BLOCK-----'
  end
end