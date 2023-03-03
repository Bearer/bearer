const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

const nunjucks = require('nunjucks');

app.get("/bad", (req, _res) => {
  nunjucks.render(req.params.body);
  nunjucks.renderString(req.params.body);
})
