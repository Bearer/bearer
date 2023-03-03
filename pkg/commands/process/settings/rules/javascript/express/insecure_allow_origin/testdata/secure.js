var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.get("/secure", (req, res) => {
  var origin = "https://some-origin"
  res.writeHead(200, { "Access-Control-Allow-Origin": "https://mish.bear" })
  res.set("access-control-allow-origin", origin)
})
