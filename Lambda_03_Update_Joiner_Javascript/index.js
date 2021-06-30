const AWS = require("aws-sdk");
const dynamo = new AWS.DynamoDB.DocumentClient();

exports.handler = async (event) => {

    console.log("Joiner: " + JSON.stringify(event, null, 4));

    console.log(" ------ Process Document [START]:  " + new Date().toLocaleTimeString());
    
    var joiner = JSON.parse(event.body);
    
    await dynamo.put({
        TableName: "Joiner",
        Item: joiner
    }).promise();
    
    console.log("Joiner: " + JSON.stringify(joiner, null, 4));
    
    console.log(" ------ Process Document [END]:    " + new Date().toLocaleTimeString());
    
    return joiner;
};