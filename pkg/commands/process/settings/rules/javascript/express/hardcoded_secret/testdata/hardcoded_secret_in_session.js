import { session } from "express-session"
var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app = express.app()

app.use(
  session({
    name: "my-custom-session-name",
    secret: "my-hardcoded-secret",
  })
)

var sessionConfig = {
  name: "my-custom-session-name",
  secret: "my-hardcoded-secret",
}

app.use(session(sessionConfig))
