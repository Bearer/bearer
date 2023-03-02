import { cookieSession } from "cookie-session";

const express = require("express");
const app = express();

app.use(cookieSession({
    domain: "example.com",
    httpOnly: true,
    secure: true,
    name: "my-custom-cookie-name",
    maxAge: 24 * 60 * 60 * 1000,
    path: "/some-path"
  })
);
