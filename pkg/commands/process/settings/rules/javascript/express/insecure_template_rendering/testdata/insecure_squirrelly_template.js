const express = require("express");
const app = express();

var Sqrl = require("squirrelly");

app.get("/bad", (req, _res) => {
  Sqrl.render(req.params.text, { name: "alvin" })
})
