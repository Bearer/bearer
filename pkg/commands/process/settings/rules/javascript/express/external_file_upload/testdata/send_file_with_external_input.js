var express = require('express');
var app = express();
var path = require('path');

app.get('/', function(req, res){
    var file = req.params.file

    res.sendFile(file);
    res.sendFile(path.resolve(file));
});