const session = require("express-session")
const express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.use(
  session({
    cookie: {
      domain: "example.com",
      secure: true,
      httpOnly: false,
      maxAge: 24 * 60 * 60 * 1000,
      path: "/some-path",
      name: "my-custom-cookie-name",
    },
  })
)

// some other non-express cookie usage
var config = {
  hello: "world",
  cookie: "some-encrypted-string"
}
