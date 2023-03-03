import puppeteer from "puppeteer"

var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.get("/inject", async (req, res) => {
  const browser = await puppeteer.launch()
  const page = await browser.newPage()

  var content = req.body.content
  await page.setContent(content)
  await page.goto("https://" + req.query.path)

  res.send("success")
})
