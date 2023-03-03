var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.get("/insecure", (req, res) => {
  var origin = req.params.origin
  res.writeHead(200, { "Access-Control-Allow-Origin": req.params.origin })
  res.set("access-control-allow-origin", origin)
})
