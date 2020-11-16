// Created by Jonee Ryan Ty
// Copyright ACloudApp

// forgot_password

package main

import (
	"context"
	"log"
	"regexp"
	"strings"
	"time"

	"net/http"

	acaGoConfiguration "acloudapp.org/configuration"
	acaGoMongoDBModels "acloudapp.org/databases/mongodb/models"
	acaGoMongoDBAWSUtilities "acloudapp.org/databases/mongodb/providers/aws/utilities"
	acaGoMongoDBUtilities "acloudapp.org/databases/mongodb/utilities"
	acaGoUtilities "acloudapp.org/utilities"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		return acaGoMongoDBAWSUtilities.DoResponse(response, ""), nil
	}

	acaGoUtilities.PrintMilestone(mapStore, "do init")

	t := mapStore["t"].(map[string]interface{})

	// parameters
	version, _ := t["version"].(string)
	applicationType, _ := t["application_type"].(string) // ios, android, react-native-ios, react-native-android, web-js, xamarin-wm10, ionic-tizen etc
	securityHash, _ := t["security_hash"].(string)       // security_hash - a formula of version, app_type, app_secret
	log.Println("version, applicationType, securityHash:", version, applicationType, acaGoUtilities.MaskString(securityHash, "*"))

	username, _ := t["username"].(string)
	email, _ := t["email"].(string)
	// log.Println("username, email:", username, email)

	// validate security hash TODO

	// validate parameters
	formErrors := make(map[string]interface{})
	hasError = false

	// validate username
	username = strings.TrimLeft(strings.TrimSpace(username), "@") // remove trailing @
	usernameErrors := make([](map[string]interface{}), 0)

	if username != "" {
		if len(username) < acaGoConfiguration.USERNAME_MIN_LENGTH { // min length
			usernameErrors = append(usernameErrors, map[string]interface{}{"message_key": "ERROR_USERNAME_MIN_LENGTH", "message_parameters": map[string]int{"USERNAME_MIN_LENGTH": acaGoConfiguration.USERNAME_MIN_LENGTH}})
			hasError = true
		} else {
			// regex
			usernameRe := regexp.MustCompile(acaGoConfiguration.REGEX_USERNAME_SERVER)
			if !usernameRe.MatchString(username) {
				usernameErrors = append(usernameErrors, map[string]interface{}{"message_key": "ERROR_USERNAME_REGEX"})
				hasError = true
			}
		}
	}

	if len(usernameErrors) > 0 {
		formErrors["username"] = usernameErrors
	}

	// validate email
	email = strings.TrimSpace(email)
	// email = strings.Replace(email, " ", "+", -1) // since openwhisk changes + to space for parameters, we assume that a space in the middle is a +

	emailErrors := make([]map[string]interface{}, 0)

	if email != "" {
		// regex
		emailRe := regexp.MustCompile(acaGoConfiguration.REGEX_EMAIL)
		if !emailRe.MatchString(email) {
			emailErrors = append(emailErrors, map[string]interface{}{"message_key": "ERROR_EMAIL_REGEX"})
			hasError = true
		}
	}

	if len(emailErrors) > 0 {
		formErrors["email"] = emailErrors
	}

	if username == "" && email == "" {
		usernameErrors = append(usernameErrors, map[string]interface{}{"message_key": "ERROR_USERNAME_OR_EMAIL_REQUIRED"})
		formErrors["username"] = usernameErrors
		hasError = true
	}

	if hasError {
		return acaGoMongoDBAWSUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusBadRequest,
				"body":       acaGoUtilities.GetFormErrorsJsonStringBodyReturn(false, "ERROR_VALIDATION", formErrors),
			},
			"",
		), nil
	}

	acaGoUtilities.PrintMilestone(mapStore, "parameters")

	// should be done with simple validation still need to check db

	acaGoMongoDBUtilities.DoDBConnect(mapStore)

	acaGoUtilities.PrintMilestone(mapStore, "db")

	mongoClient = mapStore["mongoClient"].(*mongo.Client) // retrieve reference

	userCol := mapStore["userCol"].(*mongo.Collection)
	// loginLogCol := mapStore["loginLogCol"].(*mongo.Collection)

	var userObj acaGoMongoDBModels.User
	// foundUserObj := false

	// find username or email if already existing
	var filter bson.M
	if username != "" && email != "" {
		filter = bson.M{"$or": []bson.M{bson.M{"username": username}, bson.M{"email": email}}}
	} else if username != "" {
		filter = bson.M{"username": username}
	} else if email != "" {
		filter = bson.M{"email": email}
	}

	findOptions := options.Find()
	ctx := context.TODO()

	var results []acaGoMongoDBModels.User
	cur, err := userCol.Find(ctx, filter, findOptions)
	cur.All(ctx, &results)

	if err == nil { // object(s) found
		if len(results) > 1 {
			return acaGoMongoDBAWSUtilities.DoResponse(
				map[string]interface{}{
					"statusCode": http.StatusBadRequest,
					"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_MULTIPLE_RESULTS", nil),
				},
				"",
			), nil

		} else if len(results) == 0 {
			log.Println("hmm possible error 1")

			return acaGoMongoDBAWSUtilities.DoResponse(
				map[string]interface{}{
					"statusCode": http.StatusNotFound,
					"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_USER_NOT_FOUND", nil),
				},
				"",
			), nil

		} else { // len == 1
			userObj = results[0]
			// foundUserObj = true
		}

	} else {
		log.Println("ERROR", err)

		return acaGoMongoDBAWSUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusNotFound,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_USER_NOT_FOUND", nil),
			},
			"",
		), nil
	}

	acaGoUtilities.PrintMilestone(mapStore, "dbs search")

	// all validations should now be done finally do save
	temporaryPassword := acaGoUtilities.GetRandomString(8, "")
	if userObj.PasswordSalt == "" {
		userObj.PasswordSalt = acaGoUtilities.GetRandomString(8, "")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(acaGoConfiguration.STATIC_SALT+userObj.PasswordSalt+temporaryPassword), bcrypt.DefaultCost)
	userObj.PasswordTemporaryHash = string(hashedPassword)
	if userObj.PasswordHash == "" {
		userObj.PasswordHash = userObj.PasswordTemporaryHash
	}

	userObj.PasswordTemporaryExpiry = time.Now().Add(time.Hour * acaGoConfiguration.PASSWORD_TEMPORARY_EXPIRY)
	userObj.UpdatedAt = time.Now()

	acaGoUtilities.PrintMilestone(mapStore, "preparation before save (bcrypt if any)")

	_, err = userObj.Save(mapStore)
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

	acaGoUtilities.PrintMilestone(mapStore, "db saving")

	// mongodb save should have worked - sending email for temporary password now
	if userObj.Email != "" {
		log.Println("sending email")

		if userObj.Email != "" {
			acaGoMongoDBAWSUtilities.DoEmail(mapStore, userObj, "FORGOT_PASSWORD", temporaryPassword)
		}

	} else {
		log.Println("NOT sending email")

		return acaGoMongoDBAWSUtilities.DoResponse(
			map[string]interface{}{
				"statusCode": http.StatusBadRequest,
				"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_USER_NO_EMAIL", nil),
			},
			"",
		), nil
	}

	acaGoUtilities.PrintMilestone(mapStore, "email")

	log.Println("function exec", (time.Now().UnixNano()-t1.UnixNano())/int64(time.Millisecond))

	// return success
	return acaGoMongoDBAWSUtilities.DoResponse(
		map[string]interface{}{
			"statusCode": http.StatusOK,
			"body":       acaGoUtilities.GetSimpleJsonStringBodyReturn(true, "FORGOT_PASSWORD_SUCCESS", nil),
		},
		"",
	), nil

}

func main() {
	lambda.Start(Handler)
}
