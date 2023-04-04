User.find_by(sanitize_sql(params[:oops]))
User.find_by!("oops #{params[:oops].to_i}")
User.find_by_sql("oops #{params[:oops].to_f}")
User.find_sole_by("oops #{ActiveRecord::Base.connection.quote(params[:oops])}")
User.find_sole_by(connection.quote("ok #{params[:ok]}"))

joins("INNER JOIN t_#{params[:oops].to_i}")

ActiveRecord::Base.connection.exec_query("SELECT #{connection.quote(params[:oops])}")

connection.select_all("SELECT #{connection.quote(params[:oops])}")
