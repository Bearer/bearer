class ApplicationController < ActionController::Base
  before_bugsnag_notify :add_user_info_to_bugsnag

  # Your controller code here

  private

  def add_user_info_to_bugsnag(event)
    event.set_user("9000", "bugs.nag@bugsnag.com", "Bugs Nag")
  end
end