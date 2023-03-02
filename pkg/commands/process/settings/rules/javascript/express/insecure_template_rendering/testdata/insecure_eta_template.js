const express = require("express");
const app = express();

import * as Eta from "eta";

app.get("/bad", (_req, _res) => {
  Eta.render(req.params, { name: "insecure" })
})