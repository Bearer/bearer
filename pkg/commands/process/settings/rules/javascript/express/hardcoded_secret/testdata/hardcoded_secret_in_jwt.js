import { express } from "express"
import { expressjwt } from "express-jwt";

const jwt = expressjwt

app = express.app();

app.get("/bad", expressjwt({ secret: "my-hardcoded-secret" }), function(_req, res) {
  res.sendStatus(200)
})

var secret = "my-hardcoded-secret"

app.get("/bad-2", jwt({ secret: secret }), function(_req, res) {
  res.sendStatus(200)
})
