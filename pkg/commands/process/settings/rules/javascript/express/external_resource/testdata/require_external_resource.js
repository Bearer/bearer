const express = require("express");
const app = express();

app.get("/bad", (req, _res) => {
  try {
    require(req.query.user.path)
  } catch (err) {
    // handle error
  }
})