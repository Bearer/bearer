const express = require("express");
const app = express();

import * as Handlebars from "handlebars";

app.get("/bad", (req, _res) => {
  Handlebars.precompile(req.body.user, options)
  Handlebars.compile(req.body.user, options)
})