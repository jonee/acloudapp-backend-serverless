// Created by Jonee Ryan Ty
// Copyright ACloudApp

// logout

package main

import (
	"log"
	"time"

	"net/http"

	// acaGoConfiguration "acloudapp.org/configuration"
	acaGoMongoDBModels "acloudapp.org/databases/mongodb/models"
	acaGoMongoDBAWSUtilities "acloudapp.org/databases/mongodb/providers/aws/utilities"
	acaGoMongoDBUtilities "acloudapp.org/databases/mongodb/utilities"
	acaGoUtilities "acloudapp.org/utilities"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
// type Response events.APIGatewayProxyResponse

var t0 = time.Now()
var mongoClient *mongo.Client

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// AclUser / token is required

	t1 := time.Now()
	log.Println("0-1 lambda life", (t1.UnixNano()-t0.UnixNano())/int64(time.Millisecond))

	mapStore := make(map[string]interface{})
	mapStore["last_milestone_time"] = t1

	// var err error
	hasError := false

	hasError, response := acaGoMongoDBAWSUtilities.DoInit(mapStore, request, false)
	if hasError {
		return acaGoMongoDBAWSUtilities.DoResponse(response, ""), nil
	}

	acaGoUtilities.PrintMilestone(mapStore, "do init")

	acaGoMongoDBUtilities.DoDBConnect(mapStore)

	acaGoUtilities.PrintMilestone(mapStore, "db")

	mongoClient = mapStore["mongoClient"].(*mongo.Client) // retrieve reference

	hasError, response = acaGoMongoDBAWSUtilities.DoCheckJwt(mapStore)
	if hasError {
		return acaGoMongoDBAWSUtilities.DoResponse(response, ""), nil
	}

	acaGoUtilities.PrintMilestone(mapStore, "jwt")

	// userObj := mapStore["userObj"].(acaGoMongoDBModels.User)
	loginLogObj := mapStore["loginLogObj"].(acaGoMongoDBModels.LoginLog)

	ret := make(map[string]interface{})
	success := false

	// invalidate the loginLog object
	if loginLogObj.InvalidReason == "" {
		loginLogObj.IsValid = false
		loginLogObj.InvalidReason = "LOGGED_OUT"
		loginLogObj.LoggedOutAt = time.Now()

		_, err := loginLogObj.Save(mapStore)
		if err != nil {
			log.Println("ERROR", err)

			log.Println("likely database down")

			return acaGoMongoDBAWSUtilities.DoResponse(
				map[string]interface{}{
					"statusCode": http.StatusInternalServerError,
					"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_INTERNAL", nil),
				},
				"",
			), nil
		}

		success = true
	}

	log.Println("success:", success)

	acaGoUtilities.PrintMilestone(mapStore, "loginlog save")

	log.Println("function exec", (time.Now().UnixNano()-t1.UnixNano())/int64(time.Millisecond))

	// return success
	return acaGoMongoDBAWSUtilities.DoResponse(
		map[string]interface{}{
			"statusCode": http.StatusOK,
			"body":       acaGoUtilities.GetJsonStringBodyReturn(success, ret),
		},
		"",
	), nil

}

func main() {
	lambda.Start(Handler)
}
