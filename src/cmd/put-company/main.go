package main

import (
	"api/internal/database"
	"api/internal/types"
	"api/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response represents an API Gateway response
type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

// Handler AWS Lambda
func Handler(ctx context.Context, request Request) (Response, error) {
	var companyID = request.PathParameters["companyID"]

	db, err := database.GetConnection()
	defer db.Close()
	if err != nil {
		return Response(utils.APIGateway500(err)), nil
	}

	var company types.Company
	err = json.Unmarshal([]byte(request.Body), &company)

	if err != nil {
		log.Println(err)
		return Response(utils.APIGateway500(errors.New("error unmarshalling"))), nil
	}

	err = database.UpdateCompany(db, companyID, company)

	if err != nil {
		log.Println(err)
		return Response(utils.APIGateway500(errors.New("db error"))), nil
	}

	return Response(utils.APIGateway204()), nil
}

func main() {
	lambda.Start(Handler)
}
