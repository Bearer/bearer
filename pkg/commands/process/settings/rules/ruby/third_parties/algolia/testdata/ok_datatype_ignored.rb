client = Algolia::Search::Client.create('YourApplicationID', 'YourWriteAPIKey')
index = client.init_index("my_index")

index.save_object({ user_id: user.user_id })

index.save_objects([{ user_id: user.user_id }])
