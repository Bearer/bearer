# detected hash assigment
session[:current_user] = user

user = {first_name: "Mark", last_name: "Addams"}

# supported hash assigment
# it is different type since hash properties have quotes around it
session[:current_user] = admin
admin = {"first_name": "Anna", "last_name": "Hyde"}