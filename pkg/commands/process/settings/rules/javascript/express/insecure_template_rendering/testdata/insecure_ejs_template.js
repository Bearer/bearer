const express = require("express");
const app = express();
const ejs = require('ejs');

app.get("/bad", (req, _res) => {
  let template = ejs.compile(req.body.user, options);
  template(data);
})

app.get("/bad-2", (req, _res) => {
  ejs.render(req.params.name, data, options);
})

