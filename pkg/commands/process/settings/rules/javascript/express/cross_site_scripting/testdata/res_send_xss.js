const express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.get("/bad", (req, res) => {
  res.send("<p>" + req.body.customer.name + "</p>")
})

app.get("/bad-2", (req, res) => {
  res.send("<p>" + req.body["user_id"] + "</p>")
})
