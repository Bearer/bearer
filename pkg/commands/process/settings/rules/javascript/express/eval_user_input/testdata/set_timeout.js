var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.post("/:id", (req, res) => {
  userInput = req.params.id
  var command = "new Function('" + userInput + "')"
  setTimeout(command)
})
