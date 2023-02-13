csv_string = CSV.generate do |csv|
	csv << ["email", "first_name", "last_name"]
	users.each do |user|
		csv << [
			user.email,
			user.first_name,
			user.last_name
		]
	end
end
