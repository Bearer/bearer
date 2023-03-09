app.use(session({}))
app.use(other())
app.use(express.static(__dirname + "/public"))
