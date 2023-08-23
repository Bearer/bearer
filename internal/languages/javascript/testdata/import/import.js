import lib, { f } from "library"
import { f as x } from "library"

lib.f()
f()
x()

const y = require("library")
y.f()
const { f } = y
f()

ignored.f()
