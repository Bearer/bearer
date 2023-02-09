client = Elasticsearch::Client.new(log: true)

user = { email: "someone@example.com" }

body = [
  { index: { _index: 'users', _id: '42' } },
  user
]

client.bulk(body: body)
