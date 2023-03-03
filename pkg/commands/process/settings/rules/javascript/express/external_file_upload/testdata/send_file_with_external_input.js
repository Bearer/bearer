var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())
var path = require("path")

app.get("/", function (req, res) {
  var file = req.params.file

  res.sendFile(file)
  res.sendFile(path.resolve(file))
  res.sendFile(req.params.file, {}, () => {})
  res.sendFile("file.txt", { root: path.join(__dirname, req.params.root) })
})
