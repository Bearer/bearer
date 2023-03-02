const express = require("express");
const app = express();
const doT = require('dot');

app.get("/bad", (req, _res) => {
  doT.template(req.params.template)
})