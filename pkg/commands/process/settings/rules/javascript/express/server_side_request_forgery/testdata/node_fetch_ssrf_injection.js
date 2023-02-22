import express from "express";
import fetch from 'node-fetch';

const app = express()

app.get("/inject", async (req, res) => {
  response = await fetch("https://" + req.query.path);
  res.json(response.data);
});
