import os
import os as fred
from os import environ

order_service_domain = os.environ["ORDER_SERVICE_DOMAIN"]
order_service_url = fred.environ.get('ORDER_SERVICE_URL')
user_service_host = environ.pop("USER_SERVICE_HOST")
