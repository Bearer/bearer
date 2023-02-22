import express from "express";
import puppeteer from "puppeteer";

const app = express()

app.get("/inject", async (req, res) => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();

  var content = req.body.content
  await page.setContent(content);
  await page.goto("https://"+req.query.path);

  res.send("success")
});