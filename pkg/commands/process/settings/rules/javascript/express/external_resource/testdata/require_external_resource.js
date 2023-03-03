var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.get("/bad", (req, _res) => {
  try {
    require(req.query.user.path)
  } catch (err) {
    // handle error
  }
})
