import { expressjwt } from "express-jwt";

app.get(
  "/revoked",
  expressjwt({ secret: config.secret, isRevoked: this.customRevokeCall(), algorithms: ["HS256"] }),
  function (req, res) {
    if (!req.auth.admin) return res.sendStatus(401);
    res.sendStatus(200);
  }
);