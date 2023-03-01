import { cookieSession } from "cookie-session";

const express = require("express");
const app = express();

app.use(cookieSession({
    domain: "example.com",
    httpOnly: false,
    secure: true,
    maxAge: 24 * 60 * 60 * 1000,
  })
);