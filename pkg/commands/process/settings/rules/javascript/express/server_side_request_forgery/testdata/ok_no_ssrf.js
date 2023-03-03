import puppeteer from "puppeteer"
var express = require("express")
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

app.get("/safety", async (_req, res) => {
  const browser = await puppeteer.launch()
  const page = await browser.newPage()
  await page.goto("https://mish.bearer.com")

  res.send("success")
})

app.get("/safety-2", async (req, res) => {
  var token = req.user.tokens.find((token) => token.kind === "safe")
  axios.get(`https://mish.com/bears?access_token=${token.accessToken}`)
  axios.get("https://mish.com/bears?access_token=" + token.accessToken)

  res.send("success")
})
