# :error_message	The String error message.	nil
error_message = "Error for #{user.name}"
# :error_class	The String class name of the error.	"Notice"
# :backtrace	The Array backtrace of the error.	caller
# :fingerprint	The String grouping fingerprint of the exception.	nil
# :force	Always report the exception when true, even when ignored.	false
# :sync	Send data synchronously (skips the worker) when true.	false
# :tags	The String comma-separated list of tags. nil
tags = "#{current_user.first_name},#{current_user.last_name}"
# :context	The Hash context to associate with the exception.	nil
context = {
  user: {
    first_name: "first name",
    last_name: "last name"
  }
}
# :controller	The String controller name (such as a Rails controller). nil
# :action	The String action name (such as a Rails controller action).	nil
# :parameters	The Hash HTTP request paramaters.	nil
parameters = {
  user: {
    email: "user@example.com"
  }
}
# :session	The Hash HTTP request session. nil
# :url	The String HTTP request URL. nil

Honeybadger.notify(
  "Something is wrong here for " + user.gender,
  class_name: "MyError",
  error_message: error_message,
  tags: tags,
  context: context,
  parameters: parameters,
)