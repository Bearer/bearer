import { CryptoJS } from "crypto-js"

var hash1 = CryptoJS.MD5(user.email);
var hash2 = CryptoJS.DES.encrypt(customer.email, "secret key");
