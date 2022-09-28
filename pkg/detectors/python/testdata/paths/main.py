def get_only_path(platform):
  messages = app.delivery_client.call(
    "GET",
    "/api/delivery-messages"
  )

  if messages.error:
    raise RuntimeError

  messages

def get_with_url_concatenated(platform_id, customer_id):
  delivery = app.delivery_client.call(
    "GET",
    app.api_url + "/api/customers" + customer_id + "/transactions/" + platform_id
  )

  if delivery.error:
    raise RuntimeError

  delivery.id

def get_with_url_passed_as_arg(platform, customer):
  delivery = app.delivery_client.call(
    "GET",
    app.api_url,
    "/api/customers" + customer.ID + "/transactions/" + platform.foo.id
  )

  if delivery.error:
    raise RuntimeError

  delivery.id

def get_with_url_fully_interpolated(platform_id, customer_id):
  delivery = app.delivery_client.call(
    "GET",
    f"{ENV['CUSTOMERS_HOST']}:{ENV['CUSTOMER_HOST_PORT']}/api/delivery-messages?num_page={page}&filters[]={filters}",
  )

  if delivery.error:
    raise RuntimeError

  delivery.id

def get_transactions_with_url_interpolated(platform_id, customer_id):
  delivery = app.delivery_client.call(
    "GET",
    f"{app.apiURL}api/customers/" + customer_id + "/transactions/" + platform_id
  )

  if delivery.error:
    raise RuntimeError

  delivery.id

def app(self):
    """ Bottle application handling this request. """
    raise RuntimeError('This request is not connected to an application.')