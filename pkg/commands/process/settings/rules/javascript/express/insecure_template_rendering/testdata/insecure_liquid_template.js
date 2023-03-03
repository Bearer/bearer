const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())


import { Liquid } from 'liquidjs'
const engine = new Liquid()

app.get("/bad", (req, _res) => {
  engine.render(req.params.text, { hello: "world" })
})