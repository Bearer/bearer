import { cookieSession } from "cookie-session"

const express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.use(
  cookieSession({
    domain: "example.com",
    httpOnly: false,
    secure: true,
    maxAge: 24 * 60 * 60 * 1000,
  })
)
