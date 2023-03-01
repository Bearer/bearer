import { CryptoJS } from "crypto-js"

var hash = CryptoJS.RC4.encrypt(user.password, "secret key");
