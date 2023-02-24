import { DynamoDBClient } from "@aws-sdk/client-dynamodb";
import { QueryCommand } from "@aws-sdk/client-dynamodb";

const dynamodb = new DynamoDBClient({ region: "af-south-1" })

exports.handler = async function(event, _context) {
  const params = {
    KeyConditionExpression: "Title = " + getTitle(),
    FilterExpression: "contains (Author, :name)",
    ExpressionAttributeValues: {
      ":name": { S: getAuthorName() },
    },
    ProjectionExpression: "Title, Author",
    TableName: "BOOKS_TABLE",
  };

  const data = await ddbClient.send(new QueryCommand(params));
}