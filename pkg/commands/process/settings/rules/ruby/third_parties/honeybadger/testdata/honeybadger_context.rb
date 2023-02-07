tags = "#{current_user.first_name},#{current_user.last_name}"

Honeybadger.context({
  tags: tags
})

Honeybadger.context({
  my_data: current_user.email
})