const express = require("express");
const app = express();

app.get("/bad", (req, res) => {
  res.send("<p>" + req.body.customer.name + "</p>")
})

app.get("/bad-2", (req, res) => {
  res.send("<p>" + req.body["user_id"] + "</p>")
})