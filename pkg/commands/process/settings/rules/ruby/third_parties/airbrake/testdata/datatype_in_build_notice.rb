notice = Airbrake.build_notice('App crashed!')
notice[:params][:username] = user.name
airbrake.notify_sync(notice)