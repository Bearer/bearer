const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

const _ = require('lodash');

app.get("/bad", (req, _res) => {
  var compiled = _.template(req.params.body);
})
