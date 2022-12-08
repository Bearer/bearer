# simple hashes are supported
session[:current_user] = {first_name: "Mark", last_name: "Addams"}


# complex hashses are partially supported, only first level properties (first_name, last_name, address)
session[:current_user] = {first_name: "Anna", last_name: "Hyde", address: { zip: 10000, city: "Zagreb" }}


# should ignore since none are part of undefined datatype
admin.username = "admin"
user.email.domain = "gmail.com"
{"president_first_name": "George", president_last_name: "Washington"}