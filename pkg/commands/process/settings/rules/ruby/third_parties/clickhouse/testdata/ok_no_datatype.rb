Clickhouse.connection.insert_rows(events, products: %w(id year date amount)) do |rows|
  @products.each do |product|
    rows << [
      "123",
      Date.now.year,
      DateTime.now,
      product.stock_count
    ]
  end
end