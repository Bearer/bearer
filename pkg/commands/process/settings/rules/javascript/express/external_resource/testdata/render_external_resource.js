const express = require("express");
const app = express();

app.get("/bad", (req, res) => {
  return res.render(
    req.query.path + "/results",
    { page: 1 }
  )
})