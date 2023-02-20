var express    = require('express')
var app = express()

app.use('/ftp', express.static('public/ftp'))

app.listen(3000)