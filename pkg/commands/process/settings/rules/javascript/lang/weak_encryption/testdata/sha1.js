import { CryptoJS } from "crypto-js"
var crypto = require("crypto");

var key = "secret key";
var encrypted = crypto.createHmac("sha1", key).update(user.email);
var hashmd5 = crypto.createHash("sha1").update(user.email);

CryptoJS.HmacSHA1(user.email, "Key")
CryptoJS.SHA1(user.email)
