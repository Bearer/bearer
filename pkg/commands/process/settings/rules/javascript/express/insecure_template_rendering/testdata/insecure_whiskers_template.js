const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())


var whiskers = require("whiskers");

app.get("/bad", (req, _res) => {
  var context = {}
  whiskers.render(req.params.text, context)
})
