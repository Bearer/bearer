const express = require("express");
const app = express();
const _ = require('lodash');

app.get("/bad", (req, _res) => {
  var compiled = _.template(req.params.body);
})
