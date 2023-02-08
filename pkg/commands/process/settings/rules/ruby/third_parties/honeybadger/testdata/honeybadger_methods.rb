class User < ApplicationRecord
  def to_honeybadger_context
    { user: { id: id, email: email } }
  end
end