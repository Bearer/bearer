module.exports.foo = function(req, res){
  res.redirect(req.params.url)
  res.redirect(req.query.url + "/bar")
  res.redirect("https://" + req.params.url + "/bar")
  res.redirect("http://" + req.params.path + "/bar")
}
