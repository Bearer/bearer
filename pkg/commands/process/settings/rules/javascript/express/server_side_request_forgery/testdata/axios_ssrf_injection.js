import express from "express";
import axios from 'axios';

const app = express()

app.get("/inject", async (req, res) => {
  axios
    .get(req.query.path)
    .then(response => res.json(response.data))
});
