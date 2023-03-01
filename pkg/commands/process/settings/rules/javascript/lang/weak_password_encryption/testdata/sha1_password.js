var crypto = require("crypto");

var key = "secret key";
var encrypted = crypto.createHmac("sha1", key).update(user.password);
var shaHash = crypto.createHash("sha1").update(user.password);

CryptoJS.HmacSHA1(user.password, "Key")
CryptoJS.SHA1(user.password)
