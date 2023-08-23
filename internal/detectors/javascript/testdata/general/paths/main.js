function get_only_path(platform) {
  messages = await this.client.request("GET", "/api/delivery-messages")

  messages
}

function get_with_url_concatenated(platform_id, customer_id) {
  messages = await this.client.request(
    "GET",
    this.app.api_url +
      "/api/customers" +
      customer_id +
      "/transactions/" +
      platform_id
  )

  delivery.id
}

function get_with_url_passed_as_arg(platform, customer) {
  delivery = await this.client.request(
    "GET",
    this.app.api_url,
    "/api/customers" + customer.ID + "/transactions/" + platform.GetID()
  )

  delivery.id
}

function get_with_url_fully_interpolated(platform_id, customer_id) {
  delivery = await this.client.request(
    "GET",
    `${CUSTOMERS_HOST}:${CUSTOMER_HOST_PORT}/api/delivery-messages?num_page=${page}&filters[]=${filters}`
  )

  delivery.id
}

function get_transactions_with_url_interpolated(platform_id, customer_id) {
  delivery = await this.client.request(
    "GET",
    `${this.app.apiURL}api/customers/` +
      customer_id +
      "/transactions/" +
      platform_id
  )

  delivery.id
}
