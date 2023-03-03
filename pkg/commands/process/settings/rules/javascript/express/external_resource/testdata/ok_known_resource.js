var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.get("/good", (_req, res) => {
  var internalPath = "/safe-resource"
  try {
    require(internalPath)
  } catch (err) {
    // handle error
  }

  return res.render(internalPath + "/results", { page: res.params.page })
})
