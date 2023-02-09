client = Algolia::Search::Client.create('YourApplicationID', 'YourWriteAPIKey')
index = client.init_index("my_index")

index.save_object({ foo: "bar" }

index.save_objects([{ foo: "bar" }])
