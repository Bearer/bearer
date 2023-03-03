const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())


import { create } from 'express-handlebars';

const hbs = create();
app.get("/bad", (req, _res) => {
  hbs.renderView(req.params.viewPath, options, (err) => {})
})