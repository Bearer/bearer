import { CryptoJS } from "crypto-js"

var hash = CryptoJS.RC4.encrypt(user.email, "secret key");
