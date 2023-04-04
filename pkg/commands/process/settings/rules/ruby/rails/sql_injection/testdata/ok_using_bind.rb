User.find_by("attr = ?", params[:ok])
User.find_by!(attr: params[:ok])
User.find_by_sql("attr = ?", params[:ok])
User.find_sole_by(["attr = ?", params[:ok]])

ActiveRecord::Base.connection.exec_query("SELECT ?", "test", [[nil, params[:ok]]])

connection.select_all("SELECT ?", "test", [[nil, params[:ok]]])
