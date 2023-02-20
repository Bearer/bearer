var express = require('express');
var app = express();

app.get('/', function(_, res){
    res.sendFile('index.js')
});