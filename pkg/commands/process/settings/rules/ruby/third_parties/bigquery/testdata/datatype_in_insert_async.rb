bigquery = Google::Cloud::Bigquery.new
dataset = bigquery.dataset("my_dataset")

inserter = dataset.insert_async "my_table" do |result|
  call
end

inserter.insert([{ first_name: user.first_name }])
