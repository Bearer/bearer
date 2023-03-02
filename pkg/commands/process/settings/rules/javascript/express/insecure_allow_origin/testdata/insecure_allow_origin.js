var express = require('express')
var app = express()

app.get("/insecure", (req, res) => {
  var origin = req.params.origin
  res.writeHead(200, { 'Access-Control-Allow-Origin': req.params.origin })
  res.set("access-control-allow-origin", origin)
})
