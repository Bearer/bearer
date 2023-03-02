const express = require("express");
const app = express();

import { create } from 'express-handlebars';

const hbs = create();
app.get("/bad", (req, _res) => {
  hbs.renderView(req.params.viewPath, options, (err) => {})
})