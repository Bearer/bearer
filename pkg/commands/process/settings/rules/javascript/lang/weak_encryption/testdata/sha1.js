var crypto = require("crypto");

var key = "secret key";
var encrypted = crypto.createHmac("sha1", key).update(user.password);
var hashmd5 = crypto.createHash("sha1").update(user.password);
