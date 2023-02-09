client = Elasticsearch::Client.new(log: true)

user = { email: "someone@example.com" }
client.index({ index: 'users', body: user })
