const express = require("express");
const app = express();

import { Liquid } from 'liquidjs'
const engine = new Liquid()

app.get("/bad", (req, _res) => {
  engine.render(req.params.text, { hello: "world" })
})