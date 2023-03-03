const session = require("express-session")
var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.use(
  session({
    cookie: {
      domain: "example.com",
      httpOnly: true,
      secure: true,
      name: "my-custom-cookie-name",
      maxAge: 24 * 60 * 60 * 1000,
      path: "/some-path",
    },
  })
)
