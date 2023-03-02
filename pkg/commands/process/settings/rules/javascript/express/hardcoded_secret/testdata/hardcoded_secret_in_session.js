import { express } from "express"
import { session } from "express-session";

app = express.app();

app.use(session({
  name: "my-custom-session-name",
  secret: "my-hardcoded-secret"
}))

var sessionConfig = {
  name: "my-custom-session-name",
  secret: "my-hardcoded-secret"
}

app.use(session(sessionConfig))