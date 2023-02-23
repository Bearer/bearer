var express = require('express');
var app = express();

app.get('/', function(_, res){
    res.sendFile('index.js')
    res.sendFile(req.params.file, { root: path.join(__dirname, 'public') })
});
