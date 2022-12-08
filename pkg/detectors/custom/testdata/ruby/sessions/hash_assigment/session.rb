# detected hash assigment
session[:current_user] = user

user = {first_name: "Mark", last_name: "Addams"}

# not supported since admin hash proprties have quotes around it
session[:current_user] = admin
admin = {"first_name": "Anna", "last_name": "Hyde"}