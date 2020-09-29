// Created by Jonee Ryan Ty
// Copyright ACloudApp

// register

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
	// AclOpen

	t1 := time.Now()
	log.Println("0-1 lambda life", (t1.UnixNano()-t0.UnixNano())/int64(time.Millisecond))

	mapStore := make(map[string]interface{})
	mapStore["last_milestone_time"] = t1

	var err error
	hasError := false

	hasError, response := acaGoMongoDBAWSUtilities.DoInit(mapStore, request, true)
	if hasError {
		return response, nil
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
	password, _ := t["password"].(string)
	password2, _ := t["password2"].(string)
	language, _ := t["language"].(string)
	// log.Println("username, email, password, password2, lang:", username, email, acaGoUtilities.MaskString(password, "*"), acaGoUtilities.MaskString(password2, "*"), language)

	// validate security hash TODO

	// validate parameters
	formErrors := make(map[string]interface{})
	hasError = false

	// validate username
	username = strings.TrimLeft(strings.TrimSpace(username), "@") // remove trailing @
	usernameErrors := make([](map[string]interface{}), 0)

	if username == "" {
		usernameErrors = append(usernameErrors, map[string]interface{}{"message_key": "ERROR_USERNAME_REQUIRED"})
		hasError = true

	} else if len(username) < acaGoConfiguration.USERNAME_MIN_LENGTH { // min length
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

	if len(usernameErrors) > 0 {
		formErrors["username"] = usernameErrors
	}

	// validate email
	emailErrors := make([]map[string]interface{}, 0)

	if email == "" {
		emailErrors = append(emailErrors, map[string]interface{}{"message_key": "ERROR_EMAIL_REQUIRED"})
		hasError = true
	} else {
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

	// validate password
	passwordErrors := make([]map[string]interface{}, 0)

	if password == "" {
		passwordErrors = append(passwordErrors, map[string]interface{}{"message_key": "ERROR_PASSWORD_REQUIRED"})
		hasError = true
	} else if len(password) < acaGoConfiguration.PASSWORD_MIN_LENGTH {
		passwordErrors = append(passwordErrors, map[string]interface{}{"message_key": "ERROR_PASSWORD_MIN_LENGTH", "message_parameters": map[string]int{"PASSWORD_MIN_LENGTH": acaGoConfiguration.PASSWORD_MIN_LENGTH}})
		hasError = true
	}

	if len(passwordErrors) > 0 {
		formErrors["password"] = passwordErrors
	}

	// validate password2
	password2Errors := make([]map[string]interface{}, 0)

	if password != password2 {
		password2Errors = append(password2Errors, map[string]interface{}{"message_key": "ERROR_PASSWORD2_MISMATCH"})
		hasError = true
	}

	if len(password2Errors) > 0 {
		formErrors["password2"] = password2Errors
	}

	if hasError {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       acaGoUtilities.GetFormErrorsJsonStringBodyReturn(false, "ERROR_VALIDATION", formErrors),
		}, nil
	}

	if language == "" {
		language = acaGoConfiguration.DEFAULT_LANGUAGE
	}

	acaGoUtilities.PrintMilestone(mapStore, "parameters")

	// should be done with simple validation still need to check db duplicates

	acaGoMongoDBUtilities.DoDBConnect(mapStore)

	acaGoUtilities.PrintMilestone(mapStore, "db")

	mongoClient = mapStore["mongoClient"].(*mongo.Client) // retrieve reference

	userCol := mapStore["userCol"].(*mongo.Collection)
	// loginLogCol := mapStore["loginLogCol"].(*mongo.Collection)

	ctx := context.Background()

	// find username or email if already existing
	var filter bson.M
	if username != "" && email != "" {
		filter = bson.M{"$or": []bson.M{bson.M{"username": username}, bson.M{"email": email}}}
	} else if username != "" {
		filter = bson.M{"username": username}
	} else if email != "" {
		filter = bson.M{"email": email}
	}

	var tmpUserObj acaGoMongoDBModels.User
	err = userCol.FindOne(ctx, filter).Decode(&tmpUserObj)
	if err == nil { // object found so we return an error
		if tmpUserObj.Username == username {
			usernameErrors = append(usernameErrors, map[string]interface{}{"message_key": "ERROR_USERNAME_UNIQUE"})
			formErrors["username"] = usernameErrors
			hasError = true
		}
		if tmpUserObj.Email == email {
			emailErrors = append(emailErrors, map[string]interface{}{"message_key": "ERROR_EMAIL_UNIQUE"})
			formErrors["email"] = emailErrors
			hasError = true
		}

		if !hasError {
			log.Println("hmm possible error 1")
		}

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       acaGoUtilities.GetFormErrorsJsonStringBodyReturn(false, "ERROR_VALIDATION", formErrors),
		}, nil
	}

	acaGoUtilities.PrintMilestone(mapStore, "user db duplicates check")

	// all validations should now be done finally do save
	var userObj acaGoMongoDBModels.User
	userObj.Username = username
	userObj.Email = email

	// password
	userObj.PasswordSalt = acaGoUtilities.GetRandomString(8, "")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userObj.PasswordSalt+password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("ERROR", err)
	}
	userObj.PasswordHash = string(hashedPassword)

	// email validation
	userObj.IsEmailValidated = false
	userObj.ValidationSecret = acaGoUtilities.GetRandomString(8, "")

	userObj.Access = "C"

	userObj.Language = language

	userObj.CreatedAt = time.Now()

	acaGoUtilities.PrintMilestone(mapStore, "preparation before save (bcrypt if any)")

	_, err = userObj.Save(mapStore)
	if err != nil {
		log.Println("ERROR", err)

		log.Println("likely database down")

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       acaGoUtilities.GetSimpleJsonStringBodyReturn(false, "ERROR_INTERNAL", nil),
		}, nil
	}

	acaGoUtilities.PrintMilestone(mapStore, "db saving / create user")

	// mongodb save should have worked - sending email validation now
	if userObj.Email != "" {
		log.Println("sending email")

		if userObj.Email != "" {
			acaGoMongoDBAWSUtilities.DoEmail(mapStore, userObj, "VALIDATION", "")
		}

	} else {
		log.Println("NOT sending email")
	}

	acaGoUtilities.PrintMilestone(mapStore, "email")

	log.Println("function exec", (time.Now().UnixNano()-t1.UnixNano())/int64(time.Millisecond))

	// return success
	return events.APIGatewayProxyResponse{
		// IsBase64Encoded: false,
		StatusCode: http.StatusOK,
		Body:       acaGoUtilities.GetJsonStringBodyReturn(true, userObj.ExportArrayPublic()),
	}, nil

}

func main() {
	lambda.Start(Handler)
}
