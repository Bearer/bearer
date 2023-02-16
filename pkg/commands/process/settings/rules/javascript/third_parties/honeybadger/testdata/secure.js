const Honeybadger = require("@honeybadger-io/js");

let context = { user: { uuid: "aacd05fd-8f5b-4bc6-aa8b-35e5fbf37325" } };

Honeybadger.setContext(context);
