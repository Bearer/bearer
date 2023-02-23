const express = require("express");
const app = express();

var config = require(global.baseUrl + '/config.js');

app.get("/good", (_req, res) => {
  res.set('X-Frame-Options', config.options)
  res.send(200)
})