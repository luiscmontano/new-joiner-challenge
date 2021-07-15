const AWS = require("aws-sdk");
const dynamo = new AWS.DynamoDB.DocumentClient();

exports.handler = async (event) => {

    console.log(" ------ Delete Joiner [START]:  " + new Date().toLocaleTimeString());
    
    var id = parseInt(event.pathParameters.id, 10);
    
    await dynamo.delete({
        TableName: "Joiner",
        Key: {
            identificationNumber: id
        }

    }).promise();
    
    console.log(" ------ Delete Joiner [END]:    " + new Date().toLocaleTimeString());
    
    return event.pathParameters.id;
};