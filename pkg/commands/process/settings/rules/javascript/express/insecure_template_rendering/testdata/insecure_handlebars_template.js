const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())


import * as Handlebars from "handlebars";

app.get("/bad", (req, _res) => {
  Handlebars.precompile(req.body.user, options)
  Handlebars.compile(req.body.user, options)
})