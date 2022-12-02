CSV.open("path/to/user.csv", "wb") do |csv|
  csv << ["email", "first_name", "last_name"]
	users.each do |user|
		csv << [
			user.email,
			user.first_name,
			user.last_name
		]
	end
end

File.open("users.log", "w") { |f| f.write "#{Time.now} - User #{user_2.email} logged in\n" }

File.open(user_3.emails, "users.csv", "w") do |f|
	users.each do |user_4|
		f.write "#{user_4.email},#{user_4.first_name},#{user_4.last_name}"
	end
end
