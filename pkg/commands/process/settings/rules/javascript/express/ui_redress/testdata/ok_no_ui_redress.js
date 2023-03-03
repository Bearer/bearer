const express = require("express")
const app = express()
const helmet = require("helmet")

app.use(helmet())
app.disable("x-powered-by")

var config = require(global.baseUrl + "/config.js")

app.get("/good", (_req, res) => {
  res.set("X-Frame-Options", config.options)
  res.send(200)
})
