import express from "express";
import puppeteer from "puppeteer";

const app = express()

app.get("/safety", async (_req, res) => {
  const browser = await puppeteer.launch();
  const page = await browser.newPage();
  await page.goto("https://mish.bearer.com");

  res.send("success")
});