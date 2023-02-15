bigquery = Google::Cloud::Bigquery.new
dataset = bigquery.dataset("my_dataset")

inserter = dataset.insert_async "my_table" do |result|
  call
end

user = { "foo" => "bar" }
rows = [user]

inserter.insert(rows)
dataset.insert(rows)
