// main.go
/*

https://github.com/aws/aws-lambda-go

set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o main main.go
%USERPROFILE%\Go\bin\build-lambda-zip.exe -o main.zip main

https://github.com/awsdocs/aws-doc-sdk-examples/tree/master/go/example_code/dynamodb

*/

package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    //"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    //*"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"
    //"github.com/aws/aws-sdk-go/service/dynamodb/expression"
    "net/http"
    //"encoding/json"

    "log"
    //"fmt"
)

type Task struct{
    ParentTaskId string `json:"parentTaskId"`
    Id string `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
    EstimatedRequiredHours int16 `json:"estimatedRequiredHours"`
    Stack string `json:"stack"`
    MinRole []string `json:"minRole"`
    ResourceId string `json:"resourceId"`
}


func main() {

    log.Print("Hello from main")

    lambda.Start(Handler)
}


func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

    log.Print("Hello from Handler")
    log.Print(req)
    log.Print("Parameters ")
    id := req.PathParameters["id"]
    log.Print(id)
    
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-2")},
    )

    if err != nil {
        panic(err)
    }

    //proj := expression.NamesList(expression.Name("id"))
    //expr, err := expression.NewBuilder().WithProjection(proj).Build()

    svc := dynamodb.New(sess)

    //taskIDs := []string { param1 }

    //mapOfAttrKeys := []map[string]*dynamodb.AttributeValue{}

    /*for _, task := range taskIDs {
        mapOfAttrKeys = append(mapOfAttrKeys, map[string]*dynamodb.AttributeValue{
            "id": &dynamodb.AttributeValue{
                S: aws.String(task),
            },
            "attr": &dynamodb.AttributeValue{
                S: aws.String("task"),
            },
        })
    }*/

    out, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
        Key: map[string]*dynamodb.AttributeValue{
            "id": {
				S: aws.String(id),
			},
        },
        TableName: aws.String("Task"),
    })
    if err != nil {
        panic(err)
    }

    log.Print(out)

    /*var tasks = []Task{}
    var error = dynamodbattribute.UnmarshalListOfMaps(out.Responses["Task"], &tasks)
    if error != nil {
        panic (error)
    }

    fmt.Println(tasks)

    outToJson, err := json.Marshal(tasks)
    if err != nil {
        panic (err)
    }*/

    return response(string(id), http.StatusOK), nil
    //return response(out.Items, http.StatusOK), nil
}

func response(body string, statusCode int) events.APIGatewayProxyResponse {
    return events.APIGatewayProxyResponse {
        StatusCode: statusCode,
        Body: string(body),
        Headers: map[string]string {
            "Access-Control-Allow-Origin": "*",
            "Content-Type": "application/json",
        },
    }
}