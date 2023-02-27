const express = require("express");
const app = express();

app.get("/good", (_, res) => {
  return res.send("<p>hello world</p>")
})

app.get("/good-2", (req, res) => {
  return res.send(JSON.stringify({
    view: goodView,
    success: true,
    text: 'Good practice'
  }))
})

app.get("/good-3", () => {
  const results = newResult(search);

  res.send(results);
})

app.get("/good-4", () => {
  return res.send({ success: false, text: `User ${req.params.user_id} not found` });
})