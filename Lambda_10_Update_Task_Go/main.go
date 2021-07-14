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
}


func main() {

	log.Print("Hello from main")

	lambda.Start(Handler)
}


func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Print("Hello from Handler")

	var task Task

	err := json.Unmarshal([]byte(req.Body), &task)

	log.Printf(" Task In Json ")
	fmt.Printf("marshalled struct: %+v", task)

	if err != nil {
		return response("Couldn't unmarshal json into task struct", http.StatusBadRequest), nil
	}

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

    taskToJson, err := json.Marshal(task)
    if err != nil {
        fmt.Println(err)
    }

	return response(string(taskToJson), http.StatusOK), nil
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