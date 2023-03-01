import { CryptoJS } from "crypto-js"

var hash1 = CryptoJS.MD5(user.email);
var hash2 = CryptoJS.DES.encrypt(customer.email, "secret key");

app.get("/bad", async (_req, res) => {
  const hash = await argon2.hash(currentUser.email, {
    type: argon2.argon2i,
    memoryCost: 2 ** 16,
    hashLength: 50,
  });

  // do something

  return res.status(200)
})
