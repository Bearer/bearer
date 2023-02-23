var libxml = require("libxmljs");
var xml =  '<?xml version="1.0" encoding="UTF-8"?><root></root>';

libxml.parseXmlString(xml, { noent: true, noblanks: true });