var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.get("/", function (_, res) {
  res.sendFile("index.js")
  res.sendFile(req.params.file, { root: path.join(__dirname, "public") })
})
