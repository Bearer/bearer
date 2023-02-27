class MyController < ApplicationController
  skip_before_action :access_control, only: %i[public1 public2]
end
