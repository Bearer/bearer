Clickhouse.connection.insert_rows(events, users: %w(id year user_id)) do |rows|
  rows << [
    "123",
    2022,
    customer.id,
  ]
end