const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())


const Hogan = require("hogan.js");

app.get("/bad", (req, _res) => {
  var template = req.params.text
  Hogan.compile(template, { name: "insecure" })
})