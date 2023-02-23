class MyController < ApplicationController
  skip_before_action :access_control, except: :secure
end
