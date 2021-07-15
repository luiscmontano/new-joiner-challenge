const AWS = require("aws-sdk");
const dynamo = new AWS.DynamoDB.DocumentClient();

exports.handler = async(event) => {

    console.log(" ------ Get Joiners [START]:  " + new Date().toLocaleTimeString());

    var body = null;
    var ids = [];
    
    if(event.pathParameters && event.pathParameters.id){
        ids.push(parseInt(event.pathParameters.id, 10));
    } else if(event.queryStringParameters && event.queryStringParameters.identificationNumber){
        var idsFromParams = event.queryStringParameters.identificationNumber.split(",");
        for(var i = 0; i < idsFromParams.length; i++){
           ids.push(parseInt(idsFromParams[i], 10));
        }
    }
    

    if (ids.length == 0) {
        console.log("Get all items");
        body = await dynamo.scan({ TableName: "Joiner" }).promise();
    }
    else {
        console.log("Get items: " + ids);
        
        var keys = [];
        for(var i = 0; i < ids.length; i++){
            keys.push({
                identificationNumber: ids[i]
            });
        }

        var params = {
            "RequestItems": {
                "Joiner": {
                    "Keys": keys,
                }
            }
        };

        await dynamo.batchGet(params, function(err, data) {
            if (err) {
                return "Unable to read item. Error JSON:" + JSON.stringify(err, null, 2);
            }
            else {
                body = ids.length == 1 ? data.Responses.Joiner[0] : data.Responses.Joiner;
            }
        }).promise();
    }

    console.log(" ------ Get Joiners [END]:    " + new Date().toLocaleTimeString());

    return body;
};
