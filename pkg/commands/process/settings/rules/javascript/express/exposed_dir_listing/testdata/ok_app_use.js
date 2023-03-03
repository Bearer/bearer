var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.use("/ftp", express.static("public/ftp"))

app.listen(3000)
