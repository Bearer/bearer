User.find_by("attr = ?", params[:oops])
User.find_by!(attr: params[:oops])
User.find_by_sql("attr = ?", params[:oops])
User.find_sole_by(["attr = ?", params[:oops]])
