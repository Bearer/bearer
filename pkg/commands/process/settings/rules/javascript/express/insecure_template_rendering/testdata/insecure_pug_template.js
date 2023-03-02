const express = require("express");
const app = express();

var pug = require('pug');

app.get("/bad", (req, res) => {
  pug.render(req.params.name, merge(options, locals))
})

app.get("/bad-2", (req, res) => {
  pug.compile(req.params.name)
})