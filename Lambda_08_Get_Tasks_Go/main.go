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
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"net/http"
	"encoding/json"

	"log"
	"fmt"
)

type Task struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Subtasks []*Task `json:"subtasks"`
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

    // Create DynamoDB client
    svc := dynamodb.New(sess)
    // snippet-end:[dynamodb.go.create_item.session]

    av, err := dynamodbattribute.MarshalMap(task)
    if err != nil {
        log.Fatalf("Got error marshalling new task item: %s", err)
    }
    // snippet-end:[dynamodb.go.create_item.assign_struct]

 	proj := expression.NamesList(expression.Name("Id"), expression.Name("Year"), expression.Name("Rating"))
    expr, err := expression.NewBuilder().WithProjection(proj).Build()

    // snippet-start:[dynamodb.go.create_item.call]
    // Create item in table Movies
    tableName := "Task"

    input := &dynamodb.PutItemInput{
        Item:      av,
        TableName: aws.String(tableName),
    }

    _, err = svc.PutItem(input)
    if err != nil {
        log.Fatalf("Got error calling PutItem: %s", err)
    }

	return response(task.Name, http.StatusOK), nil
}

func response(body string, statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse {
		StatusCode: statusCode,
		Body: string(body),
		Headers: map[string]string {
			"Access-Control-Allow-Origin": "*",
		},
	}
}