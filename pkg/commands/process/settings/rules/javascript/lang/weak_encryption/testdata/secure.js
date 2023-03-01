var crypto = require("crypto");

var key = "secret key";
var encrypted = crypto.createHmac("sha512", key).update(user.email);
var hashmd5 = crypto.createHash("sha512").update(user.email);
