private_key = ENV.fetch("PRIVATE_JWT_KEY")
JWT.encode({ user: current_user.email }, private_key, 'HS256', {})

JWT.encode({ user: current_user.email }, ENV["SECRET_KEY"])

JWT.encode({ user_name: user.name }, Rails.application.secret_key_base)
