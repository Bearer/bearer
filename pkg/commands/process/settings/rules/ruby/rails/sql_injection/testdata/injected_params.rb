User.find_by(params[:oops])
User.find_by!("oops #{params[:oops]}")
User.find_by_sql("oops #{params[:oops]}")
User.find_sole_by("oops #{params[:oops]}")
# from
# joins
# where("str", param)
# where(["str", param])
# rewhere
# select
# select_all
# reselect
# group
# having
# order
# reorder
# delete_all
# update_all
# minimum
# maximum
# calculate($<_>, $INJECTION)
# count
# count_by_sql
# sum
# average

# chained case
# User
  # .where("oops #{params[:one]}")
  # .count("#{params[:two]}")

# ActiveRecord::Base.connection.execute
# ActiveRecord::Base.connection.exec_delete
# ActiveRecord::Base.connection.exec_insert
# ActiveRecord::Base.connection.exec_query
# ActiveRecord::Base.connection.exec_update
# connection.execute
# connection.exec_delete
# connection.exec_insert
# connection.exec_query
# connection.exec_update




# sanitize_sql
# sanitize_sql_for_assignment
# sanitize_sql_for_conditions
# sanitize_conditions
# ActiveRecord::Base.connection.quote
# connection.quote
# .to_i
# .to_f

# # aws lambda
# def handler(event:, context:)
# end

# cookies?
# request? or request.env
# params[]
