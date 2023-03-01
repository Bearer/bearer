const argon2 = require("argon2");

app.get("/bad", async (_req, res) => {
  const hash = await argon2.hash(currentUser.password, {
    type: argon2.argon2i,
    memoryCost: 2 ** 16,
    hashLength: 50,
  });

  // do something

  return res.status(200)
})

