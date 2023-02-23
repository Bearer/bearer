const express = require("express");
const app = express();

app.get("/goos", (_, res) => {
  res.send("<p>hello world</p>")
})