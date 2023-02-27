const GoogleStrategy = require("passport-google-oauth").Strategy;
const passport = require("passport");

const strategy = new GoogleStrategy({ clientSecret: "hardcodedSecret" });
passport.use(strategy);
