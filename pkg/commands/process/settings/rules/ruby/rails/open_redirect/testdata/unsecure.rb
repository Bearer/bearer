class OrdersController < ApplicationController
    def notify
      redirect_to(params[:id])
    end
  end
