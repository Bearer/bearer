const express = require("express");
const app = express();
const nunjucks = require('nunjucks');

app.get("/bad", (req, _res) => {
  nunjucks.render(req.params.body);
  nunjucks.renderString(req.params.body);
})
