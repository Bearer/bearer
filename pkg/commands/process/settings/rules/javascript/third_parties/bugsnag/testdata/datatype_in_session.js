import { Bugsnag } from "@bugsnag/js"

var bugSession = Bugsnag.startSession()
bugSession.notify(user.email)