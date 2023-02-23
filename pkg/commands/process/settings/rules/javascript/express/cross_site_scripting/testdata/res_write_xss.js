const express = require("express");
const app = express();

app.get("/bad", (req, res) => {
  var customerName = req.body.customer.name
  res.write("<h3> Greetings " + customerName + "</h3>")
})