const express = require("express");
const app = express();

var whiskers = require("whiskers");

app.get("/bad", (req, _res) => {
  var context = {}
  whiskers.render(req.params.text, context)
})
