# Not detected since it isn't supported by datatypes extraction
session[:current_user] = {first_name: Addams, last_name: "Addams"}


# should ignore since none are part of undefined datatype
admin.username = "admin"
user.email.domain = "gmail.com"
