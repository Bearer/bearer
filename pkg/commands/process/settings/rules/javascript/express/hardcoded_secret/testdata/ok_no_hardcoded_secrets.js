import { express } from "express"
import { jwt } from "express-jwt"

app = express.app();

app.get("/ok", jwt({ secret: config.get("my-secret") }), function(_req, res) {
  res.sendStatus(200)
})

var secret = process.env.SAFE_SECRET

app.get("/ok", jwt({ secret: secret }), function(_req, res) {
  res.sendStatus(200)
})

app.use(session({
  secret: config.secret,
  name: "my-custom-session-name"
}))

var sessionConfig = {
  name: "my-custom-session-name",
  secret: process.env.SAFE_SECRET
}

app.use(session(sessionConfig))