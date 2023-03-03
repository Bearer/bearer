const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

const ejs = require('ejs');

app.get("/bad", (req, _res) => {
  let template = ejs.compile(req.body.user, options);
  template(data);
})

app.get("/bad-2", (req, _res) => {
  ejs.render(req.params.name, data, options);
})

