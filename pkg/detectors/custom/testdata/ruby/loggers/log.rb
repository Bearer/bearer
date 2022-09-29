# Loggers
## Detected
log.info("Email for user: ", user.email) # 1. user 2. user.email
log.info("User email", 1, user.email, convert_contact(user)) # 1. user 2. user.email
log.info("User info:", 1, user.contact.email.address, convert_contact(user)) # 1. user 2. user.contact 3. contact.email 4. email.address
log.info(user, email)
log.info("Email for user: #{user.email}")
log.info("User name is " + user.name + " and email is " + user.email)
Rails.logger.info('whatever', address)
Rails.logger.info address2

## Not Detected
log.error(user, email)
logger.info 'your message goes here'

