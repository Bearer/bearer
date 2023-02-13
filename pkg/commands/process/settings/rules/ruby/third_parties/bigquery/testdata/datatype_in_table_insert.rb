bigquery = Google::Cloud::Bigquery.new
dataset = bigquery.dataset("my_dataset")
table = dataset.table("my_table")

user = { "first_name" => "someone" }
rows = [user]

table.insert(rows)
