// Created by Jonee Ryan Ty
// Copyright ACloudApp

// email_validation

package main

import (
	"context"
	"log"
	"strings"
	"time"

	"net/http"

	acaGoConfiguration "acloudapp.org/configuration"
	acaGoMongoDBModels "acloudapp.org/databases/mongodb/models"
	acaGoMongoDBAWSUtilities "acloudapp.org/databases/mongodb/providers/aws/utilities"
	acaGoMongoDBUtilities "acloudapp.org/databases/mongodb/utilities"
	acaGoUtilities "acloudapp.org/utilities"

	"github.com/avct/uasurfer"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	// AclOpen

	t1 := time.Now()
	log.Println("0-1 lambda life", (t1.UnixNano()-t0.UnixNano())/int64(time.Millisecond))

	mapStore := make(map[string]interface{})
	mapStore["last_milestone_time"] = t1

	var err error
	hasError := false

	hasError, response := acaGoMongoDBAWSUtilities.DoInit(mapStore, request, false)
	if hasError {
		response.Headers = map[string]string{"Content-Type": "text/html"}
		return response, nil
	}

	acaGoUtilities.PrintMilestone(mapStore, "do init")

	// t := mapStore["t"].(map[string]interface{})

	// get parameters
	queryStringParameters := request.QueryStringParameters
	language, _ := queryStringParameters["language"] // check for any in get parameters
	if language == "" {
		language = acaGoConfiguration.DEFAULT_LANGUAGE
	}
	log.Println("language:", language)

	// request path parameter
	pathParameters := request.PathParameters
	path, _ := pathParameters["path"]
	// log.Println("path:", path)

	split := strings.Split(path, "_")
	if len(split) != 2 {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-Type": "text/html"},
			StatusCode: http.StatusBadRequest,
			Body:       acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "ERROR_BAD_REQUEST", nil, language),
		}, nil
	}

	mongoIdString := acaGoUtilities.DecodeB64(split[0])
	validationSecret := split[1]

	mongoId, err := primitive.ObjectIDFromHex(mongoIdString)
	// check bson ids
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-Type": "text/html"},
			StatusCode: http.StatusBadRequest,
			Body:       acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "ERROR_BAD_REQUEST", nil, language),
		}, nil
	}

	acaGoUtilities.PrintMilestone(mapStore, "parameters")

	// should be done with simple parameter validation still need to check db

	acaGoMongoDBUtilities.DoDBConnect(mapStore)

	acaGoUtilities.PrintMilestone(mapStore, "db")

	mongoClient = mapStore["mongoClient"].(*mongo.Client) // retrieve reference

	userCol := mapStore["userCol"].(*mongo.Collection)
	// loginLogCol := mapStore["loginLogCol"].(*mongo.Collection)

	ctx := context.Background()

	// find user and validate parameters
	f := bson.M{"_id": mongoId}
	var userObj acaGoMongoDBModels.User
	err = userCol.FindOne(ctx, f).Decode(&userObj)
	if err == nil { // object found
		if userObj.ValidationSecret != validationSecret {
			return events.APIGatewayProxyResponse{
				Headers:    map[string]string{"Content-Type": "text/html"},
				StatusCode: http.StatusBadRequest,
				Body:       acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "ERROR_BAD_REQUEST", nil, language),
			}, nil

		} else if userObj.IsEmailValidated {
			return events.APIGatewayProxyResponse{
				Headers:    map[string]string{"Content-Type": "text/html"},
				StatusCode: http.StatusBadRequest,
				Body:       acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "ERROR_USER_ALREADY_VALIDATED", nil, language),
			}, nil
		}

	} else {
		// user not found
		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-Type": "text/html"},
			StatusCode: http.StatusNotFound,
			Body:       acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "ERROR_USER_NOT_FOUND", nil, language),
		}, nil
	}

	acaGoUtilities.PrintMilestone(mapStore, "db search")

	// validation should be done- email validate then save the user
	userObj.IsEmailValidated = true
	userObj.UpdatedAt = time.Now()

	uaString := mapStore["userAgent"].(string)
	// Parse() returns all attributes, including returning the full UA string last
	ua := uasurfer.Parse(uaString)
	log.Println(ua)

	_, err = userObj.Save(mapStore)
	if err != nil {
		log.Println("ERROR", err)

		log.Println("likely database down")

		return events.APIGatewayProxyResponse{
			Headers:    map[string]string{"Content-Type": "text/html"},
			StatusCode: http.StatusInternalServerError,
			Body:       acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "ERROR_INTERNAL", nil, language),
		}, nil
	}

	acaGoUtilities.PrintMilestone(mapStore, "db saving / update user")

	log.Println("function exec", (time.Now().UnixNano()-t1.UnixNano())/int64(time.Millisecond))

	body := acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "EMAIL_VALIDATION_SUCCESS_HTML", map[string]string{"BRANDING": acaGoConfiguration.BRANDING}, language)
	if ua.OS.Name == uasurfer.OSiOS {
		body = acaGoMongoDBUtilities.GetBackendTranslation(mapStore, "EMAIL_VALIDATION_SUCCESS_HTML_IOS", map[string]string{"BRANDING": acaGoConfiguration.BRANDING, "BRANDING2": acaGoConfiguration.BRANDING2}, language)
	}

	// return success
	return events.APIGatewayProxyResponse{
		// IsBase64Encoded: false,
		Headers:    map[string]string{"Content-Type": "text/html"},
		StatusCode: http.StatusOK,
		Body:       body,
	}, nil

}

func main() {
	lambda.Start(Handler)
}
