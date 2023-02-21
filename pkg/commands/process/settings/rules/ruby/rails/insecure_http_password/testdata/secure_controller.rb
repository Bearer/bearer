class UsersController < ApplicationController
    http_basic_authenticate_with name: "foo", password: Rails.secret.password
  
    def index
    end
  end