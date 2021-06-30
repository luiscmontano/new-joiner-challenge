const AWS = require("aws-sdk");
const dynamo = new AWS.DynamoDB.DocumentClient();

exports.handler = async (event) => {

    console.log(" ------ Process Document [START]:  " + new Date().toLocaleTimeString());
    
    var body = event.Records ? JSON.parse(event.Records[0].body) : event.body;
    var joiner = JSON.parse(body);

    await dynamo.put({
        TableName: "Joiner",
        Item: joiner
    }).promise();
    
    console.log("Joiner: " + JSON.stringify(joiner, null, 4));
    
    console.log(" ------ Process Document [END]:    " + new Date().toLocaleTimeString());
    
    return joiner;
};