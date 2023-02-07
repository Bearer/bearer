class ApplicationController < ActionController::Base
  before_bugsnag_notify :add_diagnostics_to_bugsnag

  # Your controller code here

  private

  def add_diagnostics_to_bugsnag(event)
    event.add_metadata(:diagnostics, {
      user: current_user.name
    })
  end
end