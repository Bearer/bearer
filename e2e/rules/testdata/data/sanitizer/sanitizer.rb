x = user.email

# unsanitized
log("abc" + user.email)
log("abc" + x)

# sanitized
y = hash(x)
log("abc" + hash(user.email))
log("abc" + y)
