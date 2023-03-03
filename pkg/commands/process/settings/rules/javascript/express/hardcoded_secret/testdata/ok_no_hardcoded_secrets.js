import { jwt } from "express-jwt"
var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app = express.app()

app.get("/ok", jwt({ secret: config.get("my-secret") }), function (_req, res) {
  res.sendStatus(200)
})

var secret = process.env.SAFE_SECRET

app.get("/ok", jwt({ secret: secret }), function (_req, res) {
  res.sendStatus(200)
})

app.use(
  session({
    secret: config.secret,
    name: "my-custom-session-name",
  })
)

var sessionConfig = {
  name: "my-custom-session-name",
  secret: process.env.SAFE_SECRET,
}

app.use(session(sessionConfig))
