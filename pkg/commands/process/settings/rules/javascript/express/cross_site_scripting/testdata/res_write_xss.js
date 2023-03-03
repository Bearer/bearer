const express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.get("/bad", (req, res) => {
  var customerName = req.body.customer.name
  res.write("<h3> Greetings " + customerName + "</h3>")
})
