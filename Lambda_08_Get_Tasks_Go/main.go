/*
Generate ZIP

set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o main main.go
%USERPROFILE%\Go\bin\build-lambda-zip.exe -o main.zip main
*/

package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-lambda-go/events"
    "net/http"
    "encoding/json"

    "log"
    "fmt"
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
    Status string `json:"status"`
}


func main() {

    log.Print("Hello from main")

    lambda.Start(Handler)
}


func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

    log.Print("Hello from Handler")
    
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-2")},
    )

    svc := dynamodb.New(sess)

    out, err := svc.Scan(&dynamodb.ScanInput{
        TableName: aws.String("Task"),
    })

    if err != nil {
        panic(err)
    }

    var tasks = []Task{}
    var error = dynamodbattribute.UnmarshalListOfMaps(out.Items, &tasks)
    if error != nil {
        panic (error)
    }

    fmt.Println(tasks)

    outToJson, err := json.Marshal(tasks)
    if err != nil {
        panic (err)
    }

    return response(string(outToJson), http.StatusOK), nil
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