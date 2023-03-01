
const argon2 = require("argon2");

app.get("/better", async (_req, res) => {
  // this is still bad for other reasons
  const hash = await argon2.hash(currentUser.password, {
    type: argon2.argon2id,
    memoryCost: 2 ** 16,
    hashLength: 50,
  });

  // do something

  return res.status(200)
})
