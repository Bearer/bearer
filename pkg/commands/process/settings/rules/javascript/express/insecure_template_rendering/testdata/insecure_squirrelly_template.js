const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())


var Sqrl = require("squirrelly");

app.get("/bad", (req, _res) => {
  Sqrl.render(req.params.text, { name: "alvin" })
})
