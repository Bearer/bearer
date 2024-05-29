def get_only_path(platform)
  messages = DeliveryClient.call(
    "GET",
    "/api/delivery-messages"
  )

  raise DeliveryNotDeliveredError unless messages

  messages
end

def get_with_url_concatenated(platform_id, customer_id)
  delivery = DeliveryClient.call(
    "GET",
    @api_url + "/api/customers" + customer_id + "/transactions/" + platform_id
  )

  raise DeliveryNotDeliveredError unless delivery

  delivery.id
end

def get_with_url_passed_as_arg(platform, customer)
  delivery = DeliveryClient.call(
    "GET",
    @api_url,
    "/api/customers" + customer.id + "/transactions/" + platform.foo.id
  )

  raise DeliveryNotDeliveredError unless delivery

  delivery.id
end

def get_with_url_fully_interpolated(platform_id, customer_id)
  delivery = DeliveryClient.call(
    "GET",
    "#{ENV['CUSTOMERS_HOST']}:#{ENV['CUSTOMER_HOST_PORT']}/api/delivery-messages?num_page=#{page}&filters[]=#{filters}",
  )

  raise DeliveryNotDeliveredError unless delivery

  delivery.id
end

def get_transactions_with_url_interpolated(platform_id, customer_id)
  delivery = DeliveryClient.call(
    "GET",
    "#{@apiURL}api/customers/" + customer_id + "/transactions/" + platform_id
  )

  raise DeliveryNotDeliveredError unless delivery

  delivery.id
end

## Ignore

def logger
  unless @logger
    if ENV["VERBOSE"]
      @logger = Logger.new(STDOUT)
    else
      # Log to file by default
      path = "/tmp/spaceship#{Time.now.to_i}_#{Process.pid}_#{Thread.current.object_id}.log"
      @logger = Logger.new(path)
    end

    @logger.formatter = proc do |severity, datetime, progname, msg|
      severity = format('%-5.5s', severity)
      "#{severity} [#{datetime.strftime('%H:%M:%S')}]: #{msg}\n"
    end
  end

  @logger
end