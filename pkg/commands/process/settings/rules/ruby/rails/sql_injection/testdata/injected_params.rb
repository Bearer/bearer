User.find_by(params[:oops])
find_by!("oops #{params[:oops]}")
User.joins("INNER JOIN #{params[:oops]}")
select("#{params[:oops]} AS oops")

# chained case
User
  .where("oops #{params[:one]}")
  .count("#{params[:two]}")

ActiveRecord::Base.connection.exec_query("SELECT #{params[:oops]}")

connection.select_all("SELECT #{params[:oops]}")
