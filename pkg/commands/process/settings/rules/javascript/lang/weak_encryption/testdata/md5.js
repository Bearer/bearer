var crypto = require("crypto");

var key = "secret key";
var encrypted = crypto.createHmac("md5", key).update(user.email);
var hashmd5 = crypto.createHash("md5").update(user.email);

import { CryptoJS } from "crypto-js"

var hash = CryptoJS.MD5(user.email);
