payload = { email: current_user.email }

JWT.encode(payload, ENV.fetch("SECRET_KEY"))
JWT.encode({ user: current_user.email }, ENV["SECRET_KEY"])

private_key = ENV.fetch("PRIVATE_JWT_KEY")
JWT.encode({ secret: "stuff", email: current_user.email }, private_key, 'HS256', {})

JWT.encode({ user_name: user.name }, Rails.application.secret_key_base)
