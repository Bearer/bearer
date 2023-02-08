class ApplicationController < ActionController::Base
  before_bugsnag_notify :add_diagnostics_to_bugsnag

  def foo
    user.email
  end

  private

  def add_diagnostics_to_bugsnag(event)
    event.add_metadata(:diagnostics, {
      user: current_user.name
    })
  end
end