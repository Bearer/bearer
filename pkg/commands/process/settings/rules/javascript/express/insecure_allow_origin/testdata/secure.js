var express = require('express')
var app = express()

app.get("/insecure", (req, res) => {
  var origin = "https://some-origin"
  res.writeHead(200, { 'Access-Control-Allow-Origin': "https://mish.bear" })
  res.set("access-control-allow-origin", origin)
})
