const Honeybadger = require("@honeybadger-io/js");

let context = { user: { email: "jhon@gmail.com" } };

Honeybadger.setContext(context);
