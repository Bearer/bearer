const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())


import * as Eta from "eta";

app.get("/bad", (_req, _res) => {
  Eta.render(req.params, { name: "insecure" })
})