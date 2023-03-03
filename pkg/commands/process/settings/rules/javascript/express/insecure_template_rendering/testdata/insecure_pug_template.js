const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())


var pug = require('pug');

app.get("/bad", (req, res) => {
  pug.render(req.params.name, merge(options, locals))
})

app.get("/bad-2", (req, res) => {
  pug.compile(req.params.name)
})