const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

const doT = require('dot');

app.get("/bad", (req, _res) => {
  doT.template(req.params.template)
})