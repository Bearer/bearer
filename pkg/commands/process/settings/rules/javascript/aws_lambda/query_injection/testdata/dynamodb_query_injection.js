// AWS SDK V3
import { DynamoDBClient } from "@aws-sdk/client-dynamodb";
import { QueryCommand } from "@aws-sdk/client-dynamodb";

const dynamodb = new DynamoDBClient({ region: "af-south-1" })

exports.handler = async function(event, _context) {
  const params = event["query_params"];

  const data = await ddbClient.send(new QueryCommand(params));
}

// AWS SDK V2
var AWS = require('aws-sdk');
var docClient = new AWS.DynamoDB.DocumentClient({apiVersion: '2012-12-20'});

exports.handler = async function(event, _context) {
  docClient.query(event["query"]["params"], function(err, data) {});
}