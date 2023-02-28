import { CryptoJS } from "crypto-js"

var hash = CryptoJS.DES.encrypt(user.email, "secret key");
