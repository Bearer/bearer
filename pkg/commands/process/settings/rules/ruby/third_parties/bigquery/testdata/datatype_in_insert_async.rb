bigquery = Google::Cloud::Bigquery.new
dataset = bigquery.dataset("my_dataset")

inserter = dataset.insert_async "my_table" do |result|
  call
end

user = { :first_name => "someone" }
rows = [user]
inserter.insert(rows)
