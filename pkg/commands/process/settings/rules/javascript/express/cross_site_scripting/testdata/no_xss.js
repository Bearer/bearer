const express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.get("/good", (_, res) => {
  return res.send("<p>hello world</p>")
})

app.get("/good-2", () => {
  // don't match on req params within strings
  return res.send({
    success: false,
    text: `User ${req.params.user_id} not found`,
  })
})

app.get("/good-3", () => {
  // don't match on custom req attributes
  const userSettings = req.user.settings
  return res.send(userSettings)
})
