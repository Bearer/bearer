var crypto = require("crypto");

var key = "secret key";
var encrypted = crypto.createHmac("md5", key).update(user.password);
var hashmd5 = crypto.createHash("md5").update(user.password);
