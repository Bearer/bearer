class MyException
  def initialize; end

  def to_airbrake
    { params: { user: current_user.email } }
  end
end

Airbrake.notify(MyException.new)