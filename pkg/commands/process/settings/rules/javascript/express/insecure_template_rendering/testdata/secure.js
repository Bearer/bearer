const express = require("express");
var helmet = require("helmet")

var app = express()
app.use(helmet())
app.use(helmet.hidePoweredBy())

// pug
var pug = require('pug');
app.get("/good-2", (_req, _res) => {
  pug.render("/mish")
})

// ejs
var ejs = require('ejs');
app.get("/good-3", (_req, _res) => {
  let template = ejs.compile(this.pageName(), options);
  template(data);
})

// lodash
const _ = require('lodash');
app.get("/good-4", (_req, _res) => {
  var compiled = _.template('<b>secure template</b>');
})

// mustache
const Mustache = require('mustache');
app.get("/good-5", (_req, _res) => {
  Mustache.render("<p>hi {{ name }} </p>", { name: "secure" })
})

// hogan.js
const Hogan = require("hogan.js");
app.get("/good-6", (_req, _res) => {
  var template = "<b>hello world</b>"
  Hogan.compile(template, { name: "secure" })
})

// squirrelly
var Sqrl = require("squirrelly");
app.get("/good-7", (_req, _res) => {
  Sqrl.render("some template", { name: "alvin" })
})

// liquid
import { Liquid } from 'liquidjs'
const engine = new Liquid()

app.get("/good-8", (_req, _res) => {
  engine.render("<h3>Hello, {{name}}!</h3>", { name: "world" })
})