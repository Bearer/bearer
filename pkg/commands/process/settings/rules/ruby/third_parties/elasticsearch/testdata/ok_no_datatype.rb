user = { foo: "bar" }
client = Elasticsearch::Client.new

body = [
  { index: { _index: 'users', _id: '42' } },
  user
]
client.bulk(body: body)

client.index(index: 'users', body: user)

client.update(index: 'books', id: 42, body: user)
