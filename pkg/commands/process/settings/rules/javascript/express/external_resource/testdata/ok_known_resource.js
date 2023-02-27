const express = require("express");
const app = express();

app.get("/good", (_req, res) => {
  var internalPath = "/safe-resource"
  try {
    require(internalPath)
  } catch (err) {
    // handle error
  }

  return res.render(
    internalPath + "/results",
    { page: res.params.page }
  )
})