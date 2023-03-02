const express = require("express");
const app = express();

const Hogan = require("hogan.js");

app.get("/bad", (req, _res) => {
  var template = req.params.text
  Hogan.compile(template, { name: "insecure" })
})