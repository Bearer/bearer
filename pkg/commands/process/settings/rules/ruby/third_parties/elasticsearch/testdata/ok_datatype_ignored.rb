user = { user_id: 42 }
client = Elasticsearch::Client.new

body = [
  { index: { _index: 'users', _id: '42' } },
  user
]
client.bulk(body: body)

client.index(index: 'users', body: user)

client.update(index: 'books', id: 42, body: user)
