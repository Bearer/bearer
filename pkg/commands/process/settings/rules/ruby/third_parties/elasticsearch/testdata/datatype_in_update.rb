user = { email: "someone@example.com" }

Elasticsearch::Client
  .new
  .update(index: 'books', id: 42, body: user)
