# trigger: data type
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

# trigger: data type
File.open("users.log", "w") { |f| f.write "#{Time.now} - User #{user_2.email} logged in\n" }

# trigger: data type
File.open(user_3.emails, "users.csv", "w") do |f|
	users.each do |user_4|
		f.write "#{user_4.email},#{user_4.first_name},#{user_4.last_name}"
	end
end

# trigger: data type
csv_string = CSV.generate do |csv|
	csv << ["email", "first_name", "last_name"]
	users.each do |user_5|
		csv << [
			user_5.email,
			user_5.first_name,
			user_5.last_name
		]
	end
end

# trigger: data type
fd = IO.sysopen("/dev/tty", "w")
IO.open(fd,"w") do |a|
  a.puts "Hello, #{user_6.full_name}!"
end
