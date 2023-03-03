const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())


const Mustache = require('mustache');

app.get("/bad", (_req, _res) => {
  Mustache.render(req.params, { name: "insecure" })
})