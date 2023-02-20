var express    = require('express')
var serveIndex = require('serve-index')
var app = express()

app.use('/public', serveIndex(__dirname + 'files'));

app.listen(3000)