import { expressjwt } from "express-jwt";

app.get(
  "/unrevoked",
  expressjwt({ secret: currentUser.email, algorithms: ["HS256"] }),
  function (req, res) {
    if (!req.auth.admin) return res.sendStatus(401);
    res.sendStatus(200);
  }
);

app.get(
  "/unrevoked",
  expressjwt({ secret: "some-secret", algorithms: ["HS256"] }),
  function (req, res) {
    if (!req.auth.admin) return res.sendStatus(401);
    res.sendStatus(200);
  }
);