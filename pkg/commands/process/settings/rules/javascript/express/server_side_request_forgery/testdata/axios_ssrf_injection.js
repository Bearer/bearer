import axios from "axios"
var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.get("/inject", async (req, res) => {
  axios.get(req.query.path).then((response) => res.json(response.data))
})
