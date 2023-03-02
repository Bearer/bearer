const express = require("express");
const app = express();

const Mustache = require('mustache');

app.get("/bad", (_req, _res) => {
  Mustache.render(req.params, { name: "insecure" })
})