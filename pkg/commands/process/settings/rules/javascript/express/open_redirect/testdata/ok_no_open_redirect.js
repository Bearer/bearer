module.exports.foo = function (_req, res) {
  res.redirect("https://google.com")
  res.redirect(!!req.query.google ? "https://google.com" : "https://bing.com")
}
