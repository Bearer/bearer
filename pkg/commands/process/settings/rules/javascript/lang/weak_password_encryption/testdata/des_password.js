import { CryptoJS } from "crypto-js"

var hash = CryptoJS.DES.encrypt(user.password, "secret key");
