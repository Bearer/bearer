# Loggers
## Detected
logger.info("Email for user: ", user.email) # 1. user 2. user.email
logger.info("User email", 1, user.email, convert_contact(user)) # 1. user 2. user.email
logger.info("User info:", 1, user.contact.email.address, convert_contact(user)) # 1. user 2. user.contact 3. contact.email 4. email.address
logger.info(user, email)
logger.info("Email for user: #{user.email}")
logger.info("User name is " + user.name + " and email is " + user.email)
Rails.logger.info('whatever', address)
Rails.logger.info address2

## Not Detected
log.error(user, email)
log.info 'your message goes here'

